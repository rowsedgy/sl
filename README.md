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




### TODO
- add edit functionality (cli and interactive)
- add connectivity status indicator
