package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/51st-state/api/pkg/rbac"
)

// CreateSchema for the cockroachdb repository
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS role_ids (
            roleId SERIAL PRIMARY KEY,
            roleIdStr TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS role_ids_idx_roleId ON role_ids (roleId);
        CREATE UNIQUE INDEX IF NOT EXISTS role_ids_idx_roleIdStr ON role_ids (roleIdStr);

        CREATE TABLE IF NOT EXISTS subject_ids (
            subjectId SERIAL PRIMARY KEY,
            subjectIdStr TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS subject_ids_idx_subjectId ON subject_ids (subjectId);
        CREATE UNIQUE INDEX IF NOT EXISTS subject_ids_idx_subjectIdStr ON subject_ids (subjectIdStr);

        CREATE TABLE IF NOT EXISTS rule_ids (
            ruleId SERIAL PRIMARY KEY,
            ruleIdStr TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS rule_ids_idx_ruleId ON rule_ids (ruleId);
        CREATE UNIQUE INDEX IF NOT EXISTS rule_ids_idx_ruleIdStr ON rule_ids (ruleIdStr);

        CREATE TABLE IF NOT EXISTS rolebindings (
            subjectId integer references subject_ids (subjectId),
            roleId integer references role_ids (roleId)
        );
        CREATE UNIQUE INDEX IF NOT EXISTS rolebindings_idx_subjectId_roleId ON rolebindings (subjectId, roleId);

        CREATE TABLE IF NOT EXISTS rulebindings (
            roleId integer references role_ids (roleId),
            ruleId integer references rule_ids (ruleId)
        );
        CREATE UNIQUE INDEX IF NOT EXISTS rulebindings_idx_roleId_ruleId ON rulebindings (roleId, ruleId);`,
	)
	return err
}

type db struct {
	database *sql.DB
}

// NewRepository for a cockroachdb and postgresql database
func NewRepository(d *sql.DB) rbac.Repository {
	return &db{
		d,
	}
}

func (d *db) GetRoleRules(ctx context.Context, roleID rbac.RoleID) (rbac.RoleRules, error) {
	rows, err := d.database.QueryContext(
		ctx,
		`SELECT rule_ids.ruleIdStr
        FROM rule_ids,
        role_ids,
        rulebindings
        WHERE role_ids.roleIdStr = $1
        AND rulebindings.roleId = role_ids.roleId
        AND rule_ids.ruleId = rulebindings.ruleId`,
		roleID,
	)
	if err != nil {
		return nil, err
	}

	roleRules := make(rbac.RoleRules, 0)
	for rows.Next() {
		var rule rbac.Rule
		if err := rows.Scan(&rule); err != nil {
			return nil, err
		}

		roleRules = append(roleRules, rule)
	}

	return roleRules, nil
}

func (d *db) upsertRoleID(ctx context.Context, roleID rbac.RoleID) error {
	_, err := d.database.ExecContext(
		ctx,
		`INSERT INTO role_ids (
            roleIdStr
        ) SELECT $1
        ON CONFLICT
        DO NOTHING`,
		roleID,
	)
	return err
}

func txError(tx *sql.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		return err
	}

	return err
}

func (d *db) SetRoleRules(ctx context.Context, roleID rbac.RoleID, rules rbac.RoleRules) error {
	if err := d.upsertRoleID(ctx, roleID); err != nil {
		return err
	}

	roleRules, err := d.GetRoleRules(ctx, roleID)
	if err != nil {
		return err
	}

	tx, err := d.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, roleRule := range roleRules {
		if !rules.Contains(roleRule) {
			if _, err := tx.ExecContext(
				ctx,
				`DELETE FROM rulebindings
                USING role_ids,
                rule_ids
                WHERE role_ids.roleIdStr = $1
                AND rule_ids.ruleIdStr = $2
                AND rulebindings.ruleId = rule_ids.ruleId
                AND rulebindings.roleId = role_ids.roleId`,
				roleID,
				roleRule,
			); err != nil {
				return txError(tx, err)
			}
		}
	}

	for _, rule := range rules {
		if !roleRules.Contains(rule) {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO rule_ids (
                    ruleIdStr
                ) SELECT $1
                ON CONFLICT
                DO NOTHING;
                INSERT INTO rulebindings (
                    roleId,
                    ruleId
                ) SELECT role_ids.roleId,
                rule_ids.ruleId
                FROM role_ids,
                rule_ids
                WHERE role_ids.roleIdStr = $2
                AND rule_ids.ruleIdStr = $3`,
				rule,
				roleID,
				rule,
			); err != nil {
				return txError(tx, err)
			}
		}
	}

	return tx.Commit()
}

func (d *db) GetSubjectRoles(ctx context.Context, subjectID rbac.SubjectID) (rbac.SubjectRoles, error) {
	rows, err := d.database.QueryContext(
		ctx,
		`SELECT role_ids.roleIdStr
        FROM role_ids,
        rolebindings,
        subject_ids
        WHERE subject_ids.subjectIdStr = $1
        AND rolebindings.subjectId = subject_ids.subjectId
        AND role_ids.roleId = rolebindings.roleId`,
		subjectID,
	)
	if err != nil {
		return nil, err
	}

	subjectRoles := make(rbac.SubjectRoles, 0)
	for rows.Next() {
		var role rbac.RoleID
		if err := rows.Scan(
			&role,
		); err != nil {
			return nil, err
		}

		subjectRoles = append(subjectRoles, role)
	}

	return subjectRoles, nil
}

func (d *db) upsertSubjectID(ctx context.Context, subjectID rbac.SubjectID) error {
	_, err := d.database.ExecContext(
		ctx,
		`INSERT INTO subject_ids (
            subjectIdStr
        ) SELECT $1
        ON CONFLICT
        DO NOTHING`,
		subjectID,
	)
	return err
}

func (d *db) SetSubjectRoles(ctx context.Context, subjectID rbac.SubjectID, roles rbac.SubjectRoles) error {
	if err := d.upsertSubjectID(ctx, subjectID); err != nil {
		return err
	}

	subjectRoles, err := d.GetSubjectRoles(ctx, subjectID)
	if err != nil {
		return err
	}

	tx, err := d.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, subjectRoleID := range subjectRoles {
		if !roles.Contains(subjectRoleID) {
			if _, err := tx.ExecContext(
				ctx,
				`DELETE FROM rolebindings
                USING role_ids,
                subject_ids
                WHERE role_ids.roleIdStr = $1
                AND subject_ids.subjectIdStr = $2
                AND rolebindings.roleId = role_ids.roleId
                AND rolebindings.subjectId = subject_ids.subjectId`,
				subjectRoleID,
				subjectID,
			); err != nil {
				return txError(tx, err)
			}
		}
	}

	for _, roleID := range roles {
		if !subjectRoles.Contains(roleID) {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO rolebindings (
                    roleId,
                    subjectId
                ) SELECT role_ids.roleId,
                subject_ids.subjectId
                FROM role_ids,
                subject_ids
                WHERE subject_ids.subjectIdStr = $1
                AND role_ids.roleIdStr = $2`,
				subjectID,
				roleID,
			); err != nil {
				return txError(tx, err)
			}
		}
	}

	return tx.Commit()
}

func (d *db) GetSubjectRuleCount(ctx context.Context, subjectID rbac.SubjectID, rule rbac.Rule) (uint64, error) {
	var count uint64
	if err := d.database.QueryRowContext(
		ctx,
		`SELECT COUNT(rulebindings.ruleId)
        FROM rolebindings,
        rulebindings,
        subject_ids,
        rule_ids
        WHERE subject_ids.subjectIdStr = $1
        AND rule_ids.ruleIdStr = $2
        AND rolebindings.subjectId = subject_ids.subjectId
        AND rulebindings.roleId = rolebindings.roleId
        AND rulebindings.ruleId = rule_ids.ruleId`,
		subjectID,
		rule,
	).Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}
