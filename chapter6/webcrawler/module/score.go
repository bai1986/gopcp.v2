package module

// CalculateScore 代表用于计算组件评分的函数类型。
type CalculateScore func(counts Counts) uint64

//代表用于计算组件评分的函数类型
type CalculateScoree func(counts Countss) uint64

// CalculateScoreSimple 代表简易的组件评分计算函数。
func CalculateScoreSimple(counts Counts) uint64 {
	return counts.CalledCount +
		counts.AcceptedCount<<1 +
		counts.CompletedCount<<2 +
		counts.HandlingNumber<<4
}

//代表简易的组件评分计算函数
func CalculateScoreSimplee(counts Countss) uint64 {
	return counts.CalledCount +
		counts.AcceptedCount<<1+
		counts.CompletedCount<<2+
		counts.HandlingNumber<<4
}

// SetScore 用于设置给定组件的评分。
// 结果值代表是否更新了评分。
func SetScore(module Module) bool {
	calculator := module.ScoreCalculator()
	//calculator是一个CalculateScore的函数类型
	if calculator == nil {
		calculator = CalculateScoreSimple
	}
	newScore := calculator(module.Counts())
	//获取当前组件评分
	if newScore == module.Score() {
		return false
	}
	//设置组件评分
	module.SetScore(newScore)
	return true
}

//setscore 用于设置给定组件的评分
//结果值代表是否更新了评分
func SetScoree(module Modulee) bool {
	calculator := module.ScoreCalculator()
	if calculator == nil {
		calculator = CalculateScoreSimplee
	}
	newScore := calculator(module.Counts())
	//获取当前组件评分
	if newScore == module.Score() {
		return false
	}
	//设置组件评分
	module.SetScore(newScore)
	return true
}
