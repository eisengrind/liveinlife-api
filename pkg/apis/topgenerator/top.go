package top

// Identifier of a top
type Identifier interface {
	Sex() uint8
	UndershirtID() uint
	UndershirtTextureID() uint
	OvershirtID() uint
	OvershirtTextureID() uint
}

type identifier struct {
	sex                                                                uint8
	undershirtID, undershirtTextureID, overshirtID, overshirtTextureID uint
}

func newIdentifier(s uint8, u, ut, o, ot uint) Identifier {
	return &identifier{s, u, ut, o, ot}
}

func (i *identifier) Sex() uint8 {
	return i.sex
}

func (i *identifier) UndershirtID() uint {
	return i.undershirtID
}

func (i *identifier) UndershirtTextureID() uint {
	return i.undershirtTextureID
}

func (i *identifier) OvershirtID() uint {
	return i.overshirtID
}

func (i *identifier) OvershirtTextureID() uint {
	return i.overshirtTextureID
}

// Provider of top data
type Provider interface {
	Data() *data
}

// Incomplete top object
type Incomplete interface {
	Provider
}

// Complete top object
type Complete interface {
	Identifier
	Incomplete
}

type data struct {
	TorsoID             uint   `json:"torso_id"`
	TorsoTextureID      uint   `json:"torso_texture_id"`
	ClothingType        uint8  `json:"clothing_type"`
	ValencyType         uint8  `json:"valency_type"`
	Name                string `json:"name"`
	Notice              string `json:"notice"`
	Status              uint8  `json:"status"`
	PolyesterPercentage uint   `json:"polyester_percentage"`
	CottonPercentage    uint   `json:"cotton_percentage"`
	LeatherPercentage   uint   `json:"leather_percentage"`
	SilkPercentage      uint   `json:"silk_percentage"`
	RelativeAmount      uint   `json:"relative_amount"`
}

func (d *data) Data() *data {
	return d
}

// SetTorsoID of a top
func (d *data) SetTorsoID(to uint) *data {
	d.TorsoID = to
	return d
}

// SetTorsoTextureID of a top
func (d *data) SetTorsoTextureID(to uint) *data {
	d.TorsoTextureID = to
	return d
}

// SetClothingType of a top
func (d *data) SetClothingType(to uint8) *data {
	d.ClothingType = to
	return d
}

// SetValencyType of a top
func (d *data) SetValencyType(to uint8) *data {
	d.ValencyType = to
	return d
}

// SetName of a top
func (d *data) SetName(to string) *data {
	d.Name = to
	return d
}

// SetNotice of a top
func (d *data) SetNotice(to string) *data {
	d.Notice = to
	return d
}

// SetStatus of a top
func (d *data) SetStatus(to uint8) *data {
	d.Status = to
	return d
}

// SetPolyesterPercentage of a top
func (d *data) SetPolyesterPercentage(to uint) *data {
	d.PolyesterPercentage = to
	return d
}

// SetCottonPercentage of a top
func (d *data) SetCottonPercentage(to uint) *data {
	d.CottonPercentage = to
	return d
}

// SetLeatherPercentage of a top
func (d *data) SetLeatherPercentage(to uint) *data {
	d.LeatherPercentage = to
	return d
}

// SetSilkPercentage of a top
func (d *data) SetSilkPercentage(to uint) *data {
	d.SilkPercentage = to
	return d
}

// SetRelativeAmount of a top
func (d *data) SetRelativeAmount(to uint) *data {
	d.RelativeAmount = to
	return d
}

// NewIncomplete top object
func NewIncomplete(name, notice string, status, clothingType, valencyType uint8, torsoID, torsoTextureID, polyester, cotton, leather, silk, relativeAmount uint) Incomplete {
	return &data{
		torsoID,
		torsoTextureID,
		clothingType,
		valencyType,
		name,
		notice,
		status,
		polyester,
		cotton,
		leather,
		silk,
		relativeAmount,
	}
}
