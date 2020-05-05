package hw03_frequency_analysis //nolint:golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Change to true if needed
var taskWithAsteriskIsCompleted = false

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		assert.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			assert.Subset(t, expected, Top10(text))
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			assert.ElementsMatch(t, expected, Top10(text))
		}
	})

	t.Run("when short text", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := "ab cd Cd cD CD 1"
			expected := []string{"cd", "ab", "1"}
			assert.Subset(t, expected, Top10(text))
		} else {
			text := "ab cd cd 1"
			expected := []string{"cd", "ab", "1"}
			assert.ElementsMatch(t, expected, Top10(text))
		}
	})

	t.Run("when only special chars in text", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `
	
   
 
			`
			var expected []string
			assert.Subset(t, expected, Top10(text))
		} else {
			text := `
	
			`
			var expected []string
			assert.ElementsMatch(t, expected, Top10(text))
		}
	})

	t.Run("when escape special chars in text", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `\n \s \t \r`
			expected := []string{"\\n", "\\s", "\\t", "\\r"}
			assert.Subset(t, expected, Top10(text))
		} else {
			text := `\n \s \t \r`
			expected := []string{"\\n", "\\s", "\\t", "\\r"}
			assert.ElementsMatch(t, expected, Top10(text))
		}
	})

	t.Run("when dash in text", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `abc - bcd-Abc bcd-aBc`
			expected := []string{"bcd-abc", "abc", "-"}
			assert.Subset(t, expected, Top10(text))
		} else {
			text := `abc - bcd-abc bcd-abc`
			expected := []string{"bcd-abc", "abc", "-"}
			assert.ElementsMatch(t, expected, Top10(text))
		}
	})
}
