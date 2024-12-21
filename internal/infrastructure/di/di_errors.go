package di

import "fmt"

func panicBindingNotFound(interfaceName string) {
	panic(fmt.Sprintf("Для интерфейса %s не зарегистрирована реализация", interfaceName))
}

func panicFirstArgumentNotInterface() {
	panic("Первый аргумент должен быть интерфейсом")
}

func panicTypeDoesntImplementInterface(typeName string, interfaceName string) {
	panic(fmt.Sprintf("Тип %s не реализует интерфейс %s", typeName, interfaceName))
}

func cantProvideInterfacePanic() {
	panic("Cant provide an interface")
}

func panicTypeCastFailed() {
	panic("Не удалось произвести приведение типов")
}

func panicDependencyNotFound(searchTypeName string) {
	panic(fmt.Sprintf("Зависимость с типом %s не найдена в контейнере зависимостей", searchTypeName))
}
