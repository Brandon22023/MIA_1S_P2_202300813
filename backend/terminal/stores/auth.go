package stores
// AuthStore es una estructura que almacena el estado de autenticación del sistema.
// Contiene información sobre si un usuario está autenticado, su nombre de usuario,
// contraseña y el ID de la partición asociada.
type AuthStore struct {
	IsLoggedIn  bool // Indica si el usuario está autenticado
	Username    string // Nombre de usuario autenticado
	Password    string // Contraseña del usuario autenticado
	PartitionID string // ID de la partición asociada al usuario autenticado
} 
// Auth es una instancia global de AuthStore que se utiliza para gestionar
// el estado de autenticación en todo el sistema.
var Auth = &AuthStore{
	IsLoggedIn:  false,
	Username:    "",
	Password:    "",
	PartitionID: "",
}

// Login establece el estado de autenticación del sistema.
// Recibe el nombre de usuario, la contraseña y el ID de la partición asociada.
// Marca al usuario como autenticado y almacena los datos proporcionados.
func (a *AuthStore) Login(username, password, partitionID string) {
	a.IsLoggedIn = true
	a.Username = username
	a.Password = password
	a.PartitionID = partitionID
}

// Logout restablece el estado de autenticación del sistema.
// Marca al usuario como no autenticado y limpia los datos almacenados.
func (a *AuthStore) Logout() {
	a.IsLoggedIn = false
	a.Username = ""
	a.Password = ""
	a.PartitionID = ""
}

// IsAuthenticated verifica si hay un usuario autenticado en el sistema.
// Devuelve true si el usuario está autenticado, de lo contrario devuelve false.
func (a *AuthStore) IsAuthenticated() bool {
	return a.IsLoggedIn
}

// GetCurrentUser devuelve los datos del usuario autenticado.
// Retorna el nombre de usuario, la contraseña y el ID de la partición asociada.
func (a *AuthStore) GetCurrentUser() (string, string, string) {
	return a.Username, a.Password, a.PartitionID
}

// GetPartitionID devuelve el ID de la partición asociada al usuario autenticado.
// Es útil para operaciones que requieren conocer la partición activa.
func (a *AuthStore) GetPartitionID() string {
	return a.PartitionID
}
