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

        CREATE TABLE IF NOT EXISTS account_ids (
            accountId SERIAL PRIMARY KEY,
            accountIdStr TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS account_ids_idx_accountId ON account_ids (accountId);
        CREATE UNIQUE INDEX IF NOT EXISTS account_ids_idx_accountIdStr ON account_ids (accountIdStr);

        CREATE TABLE IF NOT EXISTS rule_ids (
            ruleId SERIAL PRIMARY KEY,
            ruleIdStr TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS rule_ids_idx_ruleId ON rule_ids (ruleId);
        CREATE UNIQUE INDEX IF NOT EXISTS rule_ids_idx_ruleIdStr ON rule_ids (ruleIdStr);

        CREATE TABLE IF NOT EXISTS rolebindings (
            accountId integer references account_ids (accountId),
            roleId integer references role_ids (roleId)
        );
        CREATE UNIQUE INDEX IF NOT EXISTS rolebindings_idx_accountId_roleId ON rolebindings (accountId, roleId);

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

func (d *db) GetAccountRoles(ctx context.Context, accountID rbac.AccountID) (rbac.AccountRoles, error) {
	rows, err := d.database.QueryContext(
		ctx,
		`SELECT role_ids.roleIdStr
        FROM role_ids,
        rolebindings,
        account_ids
        WHERE account_ids.accountIdStr = $1
        AND rolebindings.accountId = account_ids.accountId
        AND role_ids.roleId = rolebindings.roleId`,
		accountID,
	)
	if err != nil {
		return nil, err
	}

	accountRoles := make(rbac.AccountRoles, 0)
	for rows.Next() {
		var role rbac.RoleID
		if err := rows.Scan(
			&role,
		); err != nil {
			return nil, err
		}

		accountRoles = append(accountRoles, role)
	}

	return accountRoles, nil
}

func (d *db) upsertAccountID(ctx context.Context, accountID rbac.AccountID) error {
	_, err := d.database.ExecContext(
		ctx,
		`INSERT INTO account_ids (
            accountIdStr
        ) SELECT $1
        ON CONFLICT
        DO NOTHING`,
		accountID,
	)
	return err
}

func (d *db) SetAccountRoles(ctx context.Context, accountID rbac.AccountID, roles rbac.AccountRoles) error {
	if err := d.upsertAccountID(ctx, accountID); err != nil {
		return err
	}

	accountRoles, err := d.GetAccountRoles(ctx, accountID)
	if err != nil {
		return err
	}

	tx, err := d.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, accountRoleID := range accountRoles {
		if !roles.Contains(accountRoleID) {
			if _, err := tx.ExecContext(
				ctx,
				`DELETE FROM rolebindings
                USING role_ids,
                account_ids
                WHERE role_ids.roleIdStr = $1
                AND account_ids.accountIdStr = $2
                AND rolebindings.roleId = role_ids.roleId
                AND rolebindings.accountId = account_ids.accountId`,
				accountRoleID,
				accountID,
			); err != nil {
				return txError(tx, err)
			}
		}
	}

	for _, roleID := range roles {
		if !accountRoles.Contains(roleID) {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO rolebindings (
                    roleId,
                    accountId
                ) SELECT role_ids.roleId,
                account_ids.accountId
                FROM role_ids,
                account_ids
                WHERE account_ids.accountIdStr = $1
                AND role_ids.roleIdStr = $2`,
				accountID,
				roleID,
			); err != nil {
				return txError(tx, err)
			}
		}
	}

	return tx.Commit()
}

func (d *db) GetAccountRuleCount(ctx context.Context, accountID rbac.AccountID, rule rbac.Rule) (uint64, error) {
	var count uint64
	if err := d.database.QueryRowContext(
		ctx,
		`SELECT COUNT(rulebindings.ruleId)
        FROM rolebindings,
        rulebindings,
        account_ids,
        rule_ids
        WHERE account_ids.accountIdStr = $1
        AND rule_ids.ruleIdStr = $2
        AND rolebindings.accountId = account_ids.accountId
        AND rulebindings.roleId = rolebindings.roleId
        AND rulebindings.ruleId = rule_ids.ruleId`,
		accountID,
		rule,
	).Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}
