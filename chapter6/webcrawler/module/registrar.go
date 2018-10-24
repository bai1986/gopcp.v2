package module

import (
	"fmt"
	"sync"

	"gopcp.v2/chapter6/webcrawler/errors"
)

// Registrar 代表组件注册器的接口。
type Registrar interface {
	// Register 用于注册组件实例。
	Register(module Module) (bool, error)
	// Unregister 用于注销组件实例。
	Unregister(mid MID) (bool, error)
	// Get 用于获取一个指定类型的组件的实例。
	// 本函数应该基于负载均衡策略返回实例。
	Get(moduleType Type) (Module, error)
	// GetAllByType 用于获取指定类型的所有组件实例。
	GetAllByType(moduleType Type) (map[MID]Module, error)
	// GetAll 用于获取所有组件实例。
	GetAll() map[MID]Module
	// Clear 会清除所有的组件注册记录。
	Clear()
}

//组件祖册器接口
type Registrarr interface {
	//用于注册组件实例
	Register(modulee Modulee) (bool, error)
	//取消注册组件实例
	Unregister(mid MIDD) (bool, error)
	//获取一个指定类型的组件实例
	Get(moduleType Typee) (Modulee, error)
	//获取指定类型的所有组件实例
	GetAllByType(moduleType Typee) (map[MIDD]Modulee, error)
	//虎丘所有组件实例
	GetAll() map[MIDD]Modulee
	//清楚所有组件注册记录
	Clear()
}
// NewRegistrar 用于创建一个组件注册器的实例。
func NewRegistrar() Registrar {
	return &myRegistrar{
		moduleTypeMap: map[Type]map[MID]Module{},
	}
}

//创建一个组件祖册器实例
func NewRegistrarr() Registrarr {
	return &myRegistrarr{
		moduleTypeMap:map[Typee]map[MIDD]Modulee{},
	}
}

// myRegistrar 代表组件注册器的实现类型。
type myRegistrar struct {
	// moduleTypeMap 代表组件类型与对应组件实例的映射。
	moduleTypeMap map[Type]map[MID]Module
	// rwlock 代表组件注册专用读写锁。
	rwlock sync.RWMutex
}

//组件注册器实现类型
type myRegistrarr struct {
	//代表组件类型与组件实例的映射,双重map
	moduleTypeMap map[Typee]map[MIDD]Modulee
	//代表注册专用读写锁
	rwlock sync.RWMutex
}

func (registrar *myRegistrar) Register(module Module) (bool, error) {
	if module == nil {
		return false, errors.NewIllegalParameterError("nil module instance")
	}
	mid := module.ID()
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	if !CheckType(moduleType, module) {
		errMsg := fmt.Sprintf("incorrect module type: %s", moduleType)
		return false, errors.NewIllegalParameterError(errMsg)
	}
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	modules := registrar.moduleTypeMap[moduleType]
	if modules == nil {
		modules = map[MID]Module{}
	}
	if _, ok := modules[mid]; ok {
		return false, nil
	}
	modules[mid] = module
	registrar.moduleTypeMap[moduleType] = modules
	return true, nil
}

//注册组件
func (registrarr *myRegistrarr) Register(modulee Modulee) (bool, error) {
	if modulee == nil {
		return false, errors.NewIllegalParameterErrorr("nil module instance")
	}
	mid := modulee.ID()
	parts , err := SplitMIDD(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMapp[parts[0]]
	if !CheckTypee(moduleType, modulee) {
		errMsg := fmt.Sprintf("incorrect module type: %s", moduleType)
		return false, errors.NewIllegalParameterErrorr(errMsg)
	}
	registrarr.rwlock.Lock()
	defer registrarr.rwlock.Unlock()
	modules := registrarr.moduleTypeMap[moduleType]
	if modules == nil {
		modules = map[MIDD]Modulee{}
	}
	if _, ok := modules[mid];ok {
		return false, nil
	}
	modules[mid] = modulee
	registrarr.moduleTypeMap[moduleType] = modules
	return true, nil
}

func (registrar *myRegistrar) Unregister(mid MID) (bool, error) {
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	var deleted bool
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	if modules, ok := registrar.moduleTypeMap[moduleType]; ok {
		if _, ok := modules[mid]; ok {
			delete(modules, mid)
			deleted = true
		}
	}
	return deleted, nil
}

//注销组件
func (registrar *myRegistrarr) Unregister(mid MIDD) (bool, error) {
	parts, err := SplitMIDD(mid)
	if err != nil {
		return false,err
	}
	moduleType := legalLetterTypeMapp[parts[0]]
	var deleted bool
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	if modules, ok := registrar.moduleTypeMap[moduleType]; ok {
		if _, ok := modules[mid]; ok {
			delete(modules, mid)
			deleted = true
		}
	}
	return deleted, nil
}

// Get 用于获取一个指定类型的组件的实例。
// 本函数会基于负载均衡策略返回实例。
func (registrar *myRegistrar) Get(moduleType Type) (Module, error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)
	var selectedModule Module
	for _, module := range modules {
		SetScore(module)
		if err != nil {
			return nil, err
		}
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectedModule = module
			minScore = score
		}
	}
	return selectedModule, nil
}

//get用于获取一个指定类型的组件的实例
//本函数会基于负载均衡策略返回实例
//这里会返回一个评分最低的组件
func (registrar *myRegistrarr) Get(moduleType Typee) (Modulee, error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	//查找评分最低组件
	minScore := uint64(0)
	var selectModule Modulee
	for _, module := range modules {
		SetScoree(module)
		if err != nil {
			return nil, err
		}
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectModule = module
			minScore = score
		}
	}
	return selectModule,nil
}

// GetAllByType 用于获取指定类型的所有组件实例。
func (registrar *myRegistrar) GetAllByType(moduleType Type) (map[MID]Module, error) {
	if !LegalType(moduleType) {
		errMsg := fmt.Sprintf("illegal module type: %s", moduleType)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	//允许并发读取
	registrar.rwlock.RLock()
	defer registrar.rwlock.RUnlock()
	modules := registrar.moduleTypeMap[moduleType]
	if len(modules) == 0 {
		return nil, ErrNotFoundModuleInstance
	}
	//这里为什么要重新复制一份
	//因为map属于引用类型，直接返回会有被修改的风险
	result := map[MID]Module{}
	for mid, module := range modules {
		result[mid] = module
	}
	return result, nil
}

//获取指定类型的所有组件实例
func (registrar *myRegistrarr) GetAllByType(moduleType Typee) (map[MIDD]Modulee, error) {
	if ! LegalTypee(moduleType) {
		errMsg := fmt.Sprintf("illegal module type: %s", moduleType)
		return nil, errors.NewIllegalParameterErrorr(errMsg)
	}
	//允许并发读取
	registrar.rwlock.RLock()
	defer registrar.rwlock.RUnlock()
	modules := registrar.moduleTypeMap[moduleType]
	if len(modules) == 0 {
		return nil, ErrNotFoundModuleInstance
	}
	//因为map属于引用类型，直接返回会有被修改的风险
	result := map[MIDD]Modulee{}
	for mid, module := range modules {
		result[mid] = module
	}
	return result,nil
}

// GetAll 用于获取所有组件实例。
func (registrar *myRegistrar) GetAll() map[MID]Module {
	result := map[MID]Module{}
	registrar.rwlock.RLock()
	defer registrar.rwlock.RUnlock()
	for _, modules := range registrar.moduleTypeMap {
		for mid, module := range modules {
			result[mid] = module
		}
	}
	return result
}

//获取所有组件实例
func (registrar *myRegistrarr) GetAll() map[MIDD]Modulee {
	result := map[MIDD]Modulee{}
	//允许并发读
	registrar.rwlock.RLock()
	defer registrar.rwlock.RUnlock()
	for _, modules := range registrar.moduleTypeMap {
		for mid , module := range modules {
			result[mid] = module
		}
	}
	return result
}

// Clear 会清除所有的组件注册记录。
func (registrar *myRegistrar) Clear() {
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	registrar.moduleTypeMap = map[Type]map[MID]Module{}
}

//清楚所有组件注册记录
func (registrar *myRegistrarr) Clear() {
	//清除要用写锁
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	registrar.moduleTypeMap = map[Typee]map[MIDD]Modulee{}
}

