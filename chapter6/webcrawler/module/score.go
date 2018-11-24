package module

// CalculateScore 代表用于计算组件评分的函数类型。
type CalculateScore func(counts Counts) uint64

// CalculateScoreSimple 代表简易的组件评分计算函数。
func CalculateScoreSimple(counts Counts) uint64 {
	//8<<2 == 8 * 2^2 == 8 * 4 == 32
	//每个参数所占用的权重都不一样
	//其中实时处理数据最反应一个组件是否处于忙碌中
	return counts.CalledCount +
		counts.AcceptedCount<<1 +
		counts.CompletedCount<<2 +
		counts.HandlingNumber<<4
}

// SetScore 用于设置给定组件的评分。
// 结果值代表是否更新了评分。
func SetScore(module Module) bool {
	//取出组件的评分计算器
	calculator := module.ScoreCalculator()
	//如果组件的评分计算器为空，则使用默认的评分计算器
	if calculator == nil {
		calculator = CalculateScoreSimple
	}
	//依据当前组件的各项数据统计进行实时计算新分数
	newScore := calculator(module.Counts())
	if newScore == module.Score() {
		return false
	}
	//设置组件的新分数
	module.SetScore(newScore)
	return true
}
