package di

import "reflect"

type DependencyContainer struct {
	bindings     map[reflect.Type]reflect.Type
	dependencies map[reflect.Type]reflect.Value
}

func NewDependencyContainer() *DependencyContainer {
	return &DependencyContainer{
		bindings:     make(map[reflect.Type]reflect.Type),
		dependencies: make(map[reflect.Type]reflect.Value),
	}
}

func FindDependency[T interface{}](container *DependencyContainer) *T {
	_searchType := (*T)(nil)
	searchType := reflect.TypeOf(_searchType)

	searchValue, exists := container.dependencies[searchType]
	if !exists {
		panicDependencyNotFound(searchType.Elem().Name())
	}

	convertedValue, ok := searchValue.Interface().(*T)
	if !ok {
		panicTypeCastFailed()
	}

	return convertedValue
}

func Bind[T interface{}, K interface{}](container *DependencyContainer) {
	_interfaceType := (*T)(nil)
	_implType := (*K)(nil)

	interfaceType := reflect.TypeOf(_interfaceType).Elem()
	implType := reflect.TypeOf(_implType)

	if interfaceType.Kind() != reflect.Interface {
		panicFirstArgumentNotInterface()
	}
	if !implType.Implements(interfaceType) {
		panicTypeDoesntImplementInterface(implType.Name(), interfaceType.Name())
	}

	container.bindings[interfaceType] = implType
}

func Provide[T any](container *DependencyContainer) *T {
	_object := (*T)(nil)
	objectType := reflect.TypeOf(_object).Elem()

	if objectType.Kind() == reflect.Interface {
		cantProvideInterfacePanic()
	}

	if objectValue, exists := container.dependencies[objectType]; exists {
		return castObjectType[T](objectValue)
	}

	objectValue := reflect.New(objectType)
	container.provide(objectValue.Interface())
	container.dependencies[objectType] = objectValue

	return castObjectType[T](objectValue)
}

func castObjectType[T any](objectValue reflect.Value) *T {
	castedObject, ok := objectValue.Elem().Interface().(T)
	if !ok {
		panicTypeCastFailed()
	}
	return &castedObject
}

func (c *DependencyContainer) provide(object any) {
	objectType := reflect.TypeOf(object)
	construct, exists := objectType.MethodByName("Construct")

	if !exists {
		return
	}

	constructArgs := make([]reflect.Value, 0, construct.Type.NumIn()-1)

	for i := 1; i < construct.Type.NumIn(); i++ {
		injectionType := construct.Type.In(i)

		// инъектить можно только либо указатель, либо интерфейс
		if injectionType.Kind() != reflect.Ptr && injectionType.Kind() != reflect.Interface {
			continue
		}

		if injectionType.Kind() == reflect.Interface {
			// ищем в контейнере имплементацию интерфейса
			// подставляем тип имплементации
			if injectionType, exists = c.bindings[injectionType]; !exists {
				panicBindingNotFound(injectionType.Name())
			}
		}

		if injectionValue, exists := c.dependencies[injectionType]; exists {
			constructArgs = append(constructArgs, injectionValue)
			continue
		}

		injectionValue := reflect.New(injectionType.Elem())
		c.provide(injectionValue.Interface())
		constructArgs = append(constructArgs, injectionValue)

		c.dependencies[injectionType] = injectionValue
	}

	construct.Func.Call(
		append([]reflect.Value{reflect.ValueOf(object)}, constructArgs...),
	)
}
