// Package formula 数学公式计算器
// 输入一个数学公式字符串，  例如: (3+3)*4-2
// 程序计算出答案 22, 输入的公式必须是合法的四则混合运算公式，只支持数字，以及运算符号+-*/和括号但不允许出现空格
// 使用了链表来做栈,
/* ascii 对应关系
40	(
41	)
42	*
43	+
44	,
45	-
46	.
47	/
48	0
49	1
50	2
51	3
52	4
53	5
54	6
55	7
56	8
57	9
*/

package formula

import (
	"fmt"
	"github.com/freelancer05/struct/link"
	"github.com/pkg/errors"
	"strconv"
)

// checkValid 测试字符串有效性
// 只支持数字，以及运算符号()+-*/,  括号的使用必须配对, 不允许出现空格
func checkValid(str string) error {
	err := checkChar(str)
	if err != nil {
		return err
	}

	err = checkParenthesis(str)
	if err != nil {
		return err
	}

	return nil
}

// checkChar 检查字符串是否有非法字符
func checkChar(str string) error {
	for _, s := range str {
		if s < 40 || s > 57 || s == 44 || s == 46 {
			return errors.New("Invalid string")
		}
	}
	return nil
}

// checkParenthesis 检查括号是否匹配 , ()数量需要一样
func checkParenthesis(str string) error {
	f := 0
	for _, s := range str {
		if s == 40 {
			f++
		}
		if s == 41 {
			f--
			if f < 0 {
				return errors.New("formal error")
			}
		}
	}
	if f != 0 {
		return errors.New("formal error")
	}
	return nil
}

// computer 计算器
type computer struct {
	id          int
	l           *link.SingleLink
	tmpOperator int32
}

// getValue 从链表中获取值
func (c *computer) pullValue() (string, error) {
	v, err := c.l.DelNodeByIdx(0)
	if err != nil {
		return "", err
	}
	value, ok := v.(string)
	if !ok {
		return "", err
	}
	return value, nil
}

// settle 对栈里面的数据进行计算， 遇到(则返回
func (c *computer) compute() (total float32, err error) {
	operator := ""
	for c.l.Len() > 0 {
		value, err := c.pullValue()
		if err != nil {
			return 0, err
		}

		switch value {
		case "(":
			break

		case ")":
			continue

		case "+", "-", "*", "/":
			// * + - /
			operator = value
			continue

		default:
			numberInt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, err
			}

			total, err = c.operator(operator, total, float32(numberInt))
			if err != nil {
				return 0, err
			}
			operator = ""
		}
	}

	return total, nil
}

// operator 运算
func (c *computer) operator(operator string, total, number float32) (float32, error) {
	switch operator {
	case "+":
		total += number
	case "-":
		total = number - total
	case "*":
		total *= number
	case "/":
		if total == 0 {
			return 0, errors.New("division by zero")
		}
		total = number / total
	default:
		total = number
	}
	return total, nil
}

// Controller 控制中心
type Controller struct {
	Computers map[int]*computer
	lastID    int
	tmpNumber string
}

// NewController 创建一个新的计算中心
func NewController() *Controller {
	return &Controller{Computers: make(map[int]*computer)}
}

// addComputer 添加一个新的计算器 并返回新的计算器
func (c *Controller) addComputer() *computer {
	c.lastID++
	n := &computer{l: link.NewSingleLink(), id: c.lastID}
	c.Computers[c.lastID] = n
	return n
}

// delComputer 删除一个计算器
func (c *Controller) delComputer(id int) {
	cp, ok := c.Computers[id]
	if ok {
		cp.l.Clear()
		delete(c.Computers, id)
		c.lastID--
	}
}

// getComputer 获取最后一个计算器
func (c *Controller) getComputer() *computer {
	if len(c.Computers) == 0 {
		return c.addComputer()
	}
	return c.Computers[c.lastID]
}

// isNumber 是否是数字
func (c *Controller) isNumber(s int32) bool {
	if s >= 48 || s <= 57 {
		return true
	}
	return false
}

// resolveNumb 处理数字类型
func (c *Controller) resolveNumb(s int32) error {
	c.tmpNumber += string(s)
	return nil
}

// loadingTmpNumb 将拼接的临时数字字符串填充到栈中
func (c *Controller) loadingTmpNumb() {
	if c.tmpNumber != "" {
		cp := c.getComputer()
		cp.l.InsertNode(0, c.tmpNumber)
		c.tmpNumber = ""
	}
}

// resolveOperator 处理运算符
func (c *Controller) resolveOperator(s int32) error {
	c.loadingTmpNumb()
	cp := c.getComputer()

	switch s {
	case 42:
		// *  当前一个运算符为 / 时候，需要先计算/  比如:4/2*5 期望值为10 而不是 0.4
		if cp.tmpOperator == 47 {
			r, err := cp.compute()
			if err != nil {
				return err
			}
			cp.l.InsertNode(0, fmt.Sprintf("%f", r))
		}

	case 43:
		// + 当前一个运算符为 */ 时候，需要先计算*/  比如 4*2+2  4/2+2
		if cp.tmpOperator == 42 || cp.tmpOperator == 47 {
			r, err := cp.compute()
			if err != nil {
				return err
			}
			cp.l.InsertNode(0, fmt.Sprintf("%f", r))
		}

	case 45:
		// -  当前一个运算符为 */-  时候，需要先计算*/-  比如 4*2-2  4/2-2  4-2-2
		if cp.tmpOperator == 42 || cp.tmpOperator == 47 || cp.tmpOperator == 45 {
			r, err := cp.compute()
			if err != nil {
				return err
			}
			cp.l.InsertNode(0, fmt.Sprintf("%f", r))
		}

	case 47:
		// / 当前一个运算符为 / 时候，需要先计算/  比如:8/2/2 期望值为2 而不是 8
		if cp.tmpOperator == 47 {
			r, err := cp.compute()
			if err != nil {
				return err
			}
			cp.l.InsertNode(0, fmt.Sprintf("%f", r))
		}
	}

	cp.l.InsertNode(0, string(s))
	cp.tmpOperator = s

	return nil
}

// resolveParenthesis 处理括号
func (c *Controller) resolveParenthesis(s int32) error {
	c.loadingTmpNumb()

	switch s {
	case 40:
		// (
		cp := c.addComputer()
		cp.l.InsertNode(0, string(s))

	case 41:
		// )
		cp := c.getComputer()
		cp.l.InsertNode(0, string(s))
		r, err := cp.compute()
		if err != nil {
			return err
		}
		c.delComputer(cp.id)

		newCP := c.getComputer()
		newCP.l.InsertNode(0, fmt.Sprintf("%f", r))
	}

	return nil
}

// Result 返回最终结果
func (c *Controller) Result(formulaStr string) (float32, error) {
	for _, char := range formulaStr {
		switch char {
		case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
			// 0 ~ 9
			err := c.resolveNumb(char)
			if err != nil {
				return 0, err
			}

		case 42, 43, 45, 47:
			// * + - /
			err := c.resolveOperator(char)
			if err != nil {
				return 0, err
			}

		case 40, 41:
			// ( )
			err := c.resolveParenthesis(char)
			if err != nil {
				return 0, err
			}
		}
	}

	c.loadingTmpNumb()
	cp := c.getComputer()
	return cp.compute()
}

// Calculator 四则混合运算 公式计算器
// 输入公式字符串，返回公式计算结果
func Calculator(formulaStr string) (float32, error) {
	err := checkValid(formulaStr)
	if err != nil {
		return 0, err
	}

	c := NewController()
	return c.Result(formulaStr)
}
