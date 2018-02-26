package Role

const (
	SuitDiomand = iota //方片
	SuitClub	//梅花
	SuitHeart	//红心
	SuitSpade	//黑桃
	SuitSmallKing
	SuitBigKing
)

const (
	Number0 = iota
	Number01
	Number02
	Number3
	Number4
	Number5
	Number6
	Number7
	Number8
	Number9
	Number10
	NumberJ
	NumberQ
	NumberK
	NumberA
	Number2
	NumberSmallKing
	NumberBigKing
)

type CardRole struct {
	Number uint8
	Suit   uint8
}

func (cardRole * CardRole) less(other * CardRole) bool {
	return cardRole.Number < other.Number
}

func (cardRole * CardRole) bigger(other * CardRole) bool {
	return cardRole.Number > other.Number
}

func (cardRole * CardRole) equal(other * CardRole) bool {
	return cardRole.Number == other.Number
}

func CompareCard(l, r CardRole) int {
	return (int)(l.Number - r.Number)
}

type CardSlice []CardRole

func (a CardSlice) Len() int {    	 // 重写 Len() 方法
	return len(a)
}
func (a CardSlice) Swap(i, j int){     // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a CardSlice) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
	if a[j].Number != a[i].Number {
		return a[i].Number < a[j].Number
	}
	return a[i].Suit < a[j].Suit
}