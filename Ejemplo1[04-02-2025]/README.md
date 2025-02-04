# Proyecto de Menú en Golang

## Descripción
Este proyecto es un menú interactivo en Golang que permite realizar las siguientes acciones:
1. **Ingresar un usuario** (Ejemplo de función sin parámetros).
2. **Mostrar todos los usuarios guardados** (Ejemplo de variable global).
3. **Buscar un usuario** por su nombre (Ejemplo de función con parámetro).
4. **Salir del programa**.

El proyecto está estructurado en los siguientes archivos y carpetas:
- **main.go** → Punto de entrada del programa.
- **menu/menu.go** → Contiene la lógica del menú y controla la navegación entre opciones.
- **opciones/** → Carpeta con los archivos de cada opción del menú.

---

## Conceptos Clave
### 🔹 Variables Globales
Las **variables globales** se definen fuera de cualquier función y pueden ser usadas en todo el paquete. 
Ejemplo en `opciones/opcion2.go`:

```go
var usuarios []string // Lista de usuarios (variable global)
```

### 🔹 Variables Locales
Las **variables locales** se definen dentro de una función y solo pueden ser usadas ahí.
Ejemplo en `opciones/opcion1.go`:

```go
func IngresarUsuario() {
    var nombre string // Variable local
    fmt.Print("Ingrese su nombre de usuario: ")
    fmt.Scanln(&nombre)
    usuarios = append(usuarios, nombre)
    fmt.Println("Usuario guardado correctamente.")
}
```

### 🔹 Funciones sin Parámetro
Son funciones que **no reciben valores** cuando se llaman.
Ejemplo en `opciones/opcion2.go`:

```go
func MostrarUsuarios() {
    fmt.Println("Usuarios registrados:", usuarios)
}
```

### 🔹 Funciones con Parámetro
Son funciones que **reciben valores** y los usan dentro de su lógica.
Ejemplo en `opciones/opcion3.go`:

```go
func BuscarUsuario(nombre string) {
    for _, usuario := range usuarios {
        if usuario == nombre {
            fmt.Println("El usuario", nombre, "está en el sistema.")
            return
        }
    }
    fmt.Println("El usuario", nombre, "NO está en el sistema.")
}
```

---

## Instrucciones de Uso
### 🛠️ Compilar y Ejecutar
Para compilar y ejecutar el programa, usa:
```sh
go run main.go
```

Si prefieres compilarlo y luego ejecutarlo:
```sh
go build -o menuApp main.go
./menuApp
```

### 🔹 Flujo del Menú
1️⃣ **Ingresar Usuario** → Guarda el nombre ingresado.  
2️⃣ **Mostrar Usuarios** → Muestra todos los usuarios almacenados.  
3️⃣ **Buscar Usuario** → Pide un nombre y verifica si está en la lista.  
0️⃣ **Salir** → Finaliza la ejecución.

---

## 📌 Notas Finales
- Se recomienda **usar variables locales** en lugar de globales siempre que sea posible.
- Este menú usa **un bucle infinito** para mantenerse activo hasta que el usuario seleccione "Salir".
- Se puede mejorar agregando **manejo de errores y validaciones**.

🚀 **¡Listo para probar el menú en Golang!** 😃

