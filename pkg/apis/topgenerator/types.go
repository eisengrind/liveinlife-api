package topgenerator

// Identifier of a top
type Identifier interface {
	Sex() bool
	UndershirtID() uint64
	TopID() uint64
}

type identifier struct {
	sex                 bool
	undershirtID, topID uint64
}

func newIdentifier(sex bool, undershirtID, topID uint64) Identifier {
	return &identifier{sex, undershirtID, topID}
}

func (i *identifier) Sex() bool {
	return i.sex
}

func (i *identifier) UndershirtID() uint64 {
	return i.undershirtID
}

func (i *identifier) TopID() uint64 {
	return i.topID
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
	TorsoID             uint  `json:"torso_id"`
	ClothingType        uint8 `json:"clothing_type"`
	ValencyType         uint8 `json:"valency_type"`
	Status              uint8 `json:"status"`
	PolyesterPercentage uint  `json:"polyester_percentage"`
	CottonPercentage    uint  `json:"cotton_percentage"`
	LeatherPercentage   uint  `json:"leather_percentage"`
	SilkPercentage      uint  `json:"silk_percentage"`
	RelativeAmount      uint  `json:"relative_amount"`
}

func (d *data) Data() *data {
	return d
}

// SetTorsoID of a top
func (d *data) SetTorsoID(to uint) *data {
	d.TorsoID = to
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
func NewIncomplete(status, clothingType, valencyType uint8, torsoID, polyester, cotton, leather, silk, relativeAmount uint) Incomplete {
	return &data{
		torsoID,
		clothingType,
		valencyType,
		status,
		polyester,
		cotton,
		leather,
		silk,
		relativeAmount,
	}
}
