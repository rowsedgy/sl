# SL - Herramienta de Terminal para conectividad por ssh


## Instalacion

### Binario pre-compilado

``` 
wget -O sl $(wget -q -O - https://api.github.com/repos/rowsedgy/sl/releases/latest | grep '"browser_download_url":' | grep -o 'https://[^"]*') && chmod +x sl
```


### Go install
#### Prerequisitos
- Go >= 1.25 instalado

```
go install github.com/rowsedgy/sl@latest
```
Binario se instala en **/home/\<user\>/go/bin/sl**


## Uso

El programa usa un fichero json para guardar los datos de conexiones de cada host ubicado en **/home/\<user\>/.config/sl-connections.json**.
Este fichero se puede generar a mano previamente o se generara solo al añadir algun host mediante el comando de añadir.

**\*\*** **IMPORTANTE**  **\*\***  
**TODA LA INFORMACION AÑADIDA AL FICHERO SL-CONNECTIONS ESTA GUARDADA EN TEXTO PLANO. ESTO INCLUYE CONTRASEÑAS**

El formato del archivo de conexiones es el siguiente:

```json
{
	"tunnelhosts": {
		"nombre-tunel-11": {
			"user": "Usuario SSH Tunel",
			"password": "Contraseña SSH Tunel",
			"ip": "IP tunel"
		}
	},
	"hosts": {
		"nombre-host-1": {
			"user": "Usuario SSH",
			"password": "Contraseña SSH",
			"pubauth": <true/false>,  // Define si host usa autenticacion por clave publica.
			"key": "Ruta de fichero clave publica",
			"ip": "IP de conexion por SSH",
			"webip": "URL de endpoint web", //  ej: "http://1.2.3.4/web"
			"tunnel": <true/false>, // Define si el host necesita conectividad por tunel
			"tunnelhost": "Nombre del tunel usado para conexion" // El nombre debe coincidir con el de la entrada de tunnelhosts
			"legacy": <true/false> // Define si tunel usa algoritmos antiguos ssh-rsa
		}
	}
}
```
### Modo interactivo

```
sl
```
Esto lanzara el modo interactivo. Se mostrara una lista interactiva de todos los elementos que esten dentro del fichero de conexiones.  

Keybinds basicos:
- **^ / k** - Arriba
- **v / j** - Abajo
- **q** - Salir
- **/** - Filtro
- **Enter** - Abrir conexion ssh al elemento seleccionado 
- **i** - Abrir/cerrar panel de detalles del elemento seleccionado
- **w** - Abrir enlace web del elemento seleccionado



### Modo Linea de Comandos

#### Ayuda
```
sl help
```

#### Listar todos las entradas del fichero de conexiones.  
```
sl ls
```

#### Listar todos los tuneles del fichero de conexiones.
```
sl lstun
```

#### Añadir entrada
Campos disponibles (cualquier campo no definido recibira el valor por defecto **None**)
- name 
- ip
- webip
- user
- user
- password
- pubauth
- key
- tunnel
- tunnelhost
- legcy

```
sl add --name=<nombre> --ip=<ip> --webip=<webip> --user=<usuario> --password=<contraseña> --pubauth=<true/false> --key=<ruta clave> --tunnel<true/false> --tunnelhost=<nombre tunel> --legacy=<true/false>
```

#### Añadir tunel
Campos disponibles (cualquier campo no definido recibira el valor por defecto **None**)
- name
- ip
- user
- password

```
sl addtun --name=<nombre> --ip=<ip> --user=<usuario> --password=<contraseña>
```


#### Borrar entrada
Las entradas se borran usando el campo **name**

```
sl remove --name=<nombre>
```

#### Borrar tunel
Los tuneles se borran usando el campo **name**
```
sl removetun --name=<nombre>
```

#### Editar entrada

Para editar una entrada habra que hacerlo desde el fichero de conexiones manualmente.



### Script para añadir hosts en bulk

Copiar los hosts que se quieran añadir a un fichero en formato:  
\<Nombre1\> \<IP\>  
\<Nombre2\> \<IP\>  
\<Nombre3\> \<IP\>  
...

```bash
#!/bin/bash
while read -r HOSTNAME IP _; do
    sl add --name="$HOSTNAME" --ip="$IP" --user=<usuario> --password=<contraseña>
done < $1
```

```bash
./script-bulk.sh lista-hosts.txt
```


### TODO
- add edit functionality (cli and interactive)
- add connectivity status indicator
