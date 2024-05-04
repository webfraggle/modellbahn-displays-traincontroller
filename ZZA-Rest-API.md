# Rest-API
Die Zugzielanzeiger können über eine Rest-API gesteuert werden. 
Die Übergabe von Parametern ist inkonsistent und wird je nach Funktion als GET oder POST Parameter oder als JSON Payload übergeben.

Die meisten Anfragen senden als result ein JSON.
    
    {"result":true}

Jeder Controller hat eine eigene IP-Adresse. Über einen Druck auf die Taste 0 am Controller wird die IP-Adresse angezeigt.

Um Caching zu verhindern, kann ein Random-Parameter an die URL angehängt werden z.B.:

    http://ip/next?path=GleisA&_=1714832599827 <- Random Number

Im Webinterface des Controllers kann man sich die Request im Developer-Panel unter Network anschauen.

# Ein- oder zwei Gleise
Es gibt Anzeigen mit einem oder zwei Gleisen. Die Gleise werden mit "GleisA" und "GleisB" bezeichnet.


# Nächsten oder vorherigen Zug anzeigen
Mit folgenden URLs können die Tasten für nächsten oder vorherigen Zug aus der Liste simuliert werden.

URL: **/next**

Type: **GET**

Parameter:
* path: GleisA oder GleisB

<!-- end of the list -->
    
    http://ip/next?path=GleisA

    http://ip/next?path=GleisB

    http://ip/prev?path=GleisA

    http://ip/prev?path=GleisB


# Zeit setzen
Um einen bestimmten Zug aus der Zugliste des Controllers anzuzeigen kann eine Uhrzeit an den Controller gesendet werden. Dadurch wird der Zug angezeigt, der der Zeit am nächsten ist.

URL: **/setTime**

Type: **POST**

Parameter:

- path: GleisA oder GleisB
- time: 24h-Format "12:34"

<!-- end of the list -->

# Zugtexte direkt setzen
Um die Texte der Züge direkt zu setzen muss ein JSON per POST als Payload gesendet werden. Wenn man alle drei Züge setzen will, müssen drei Request mit der jeweiligen URL gemacht werden.

URLs: 
**/zug1**
**/zug2**
**/zug3**

Type: **POST**

Payload:

    {
      "vonnach":"München",
      "nr":"RE50",
      "zeit":"00:15",
      "via":"Nürnberg",
      "abw":0,
      "hinweis":"",
      "fusszeile":"",
      "abschnitte":"",
      "reihung":"",
      "path":"GleisA"
      }

Achtung: abschnitte, fusszeile und reihung haben keine Funktion!

