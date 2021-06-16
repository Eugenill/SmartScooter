//Primero incluimos las librerías que usaremos: 
#include <PubSubClient.h> 
#include <WiFi.h>
#include <Wire.h>

// Definimos la nueva subscripción a topics
#define TEST_TOPIC "illausi@gmail.com/test"

// Replace the next variables with your Wi-Fi SSID/Password
const char *WIFI_SSID = "MiFibra-9ACF";
const char *WIFI_PASSWORD = "XYH46LjW";
char macAddress[18];    //mostra la MacAddress del Wi-Fi al que s'ha conectat

// Add MQTT Broker settings
const char *MQTT_BROKER_IP = "maqiatto.com";
const int MQTT_PORT = 1883;
const char *MQTT_USER = "illausi@gmail.com";
const char *MQTT_PASSWORD = "iria1234";
const bool RETAINED = true;
const int QoS = 0; // Quality of Service for the subscriptions

WiFiClient wifiClient;                //crea un cliente que se puede conectar a un puerto específico y una dirección IP de Internet especificada (conexión a MQTT broker)
PubSubClient mqttClient(wifiClient);  //crea una instancia de cliente parcialmente inicializada


// 2-dimensional array of row pin numbers: (pins corresponents a les files en la ESP32)
const int row[] = {15,19,32,18,13,33,12,26}; 

// 2-dimensional array of column pin numbers: (pins corresponents a les columnes en la ESP32)
const int col[] = {23,14,27,16,25,17,22,21};

//Defineixo les comandes que vull que es projectin en la matriu:
byte sp[]= {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};
byte RIGHT[]={0x00, 0x08, 0x0C, 0xFE, 0xFF, 0xFE, 0x0C, 0x08};
byte LEFT[]={0x00, 0x10, 0x30, 0x7F, 0xFF, 0x7F, 0x30, 0x10};
byte STOP[]={0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF};
byte RECTO[]={0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00};

void setup() {
  Serial.begin(9600); // Starts the serial communication

  
  for (int thisPin = 0; thisPin < 8; thisPin++)  //declara thisPin, prova que es menor que 8, incrementa thisPin en 1 (thisPin = posició que ocupa en llista ROW i COL)
  {
    // inicialitzem les sortides en els pins:
    pinMode(col[thisPin], OUTPUT);
    pinMode(row[thisPin], OUTPUT);
    //UN PIN S'ACTIVA QUAN LA ROW =low I LA COL =high
    digitalWrite(col[thisPin], LOW);
    digitalWrite(row[thisPin], HIGH);
  }
  
  
  Serial.println("\nBooting device...");

  mqttClient.setServer(MQTT_BROKER_IP, MQTT_PORT); //Es conecta la broker i a la direcció IP
  mqttClient.setCallback(callback); // funció que ens permet rebre missatges

  connectToWiFiNetwork(); // Connects to the configured network
  connectToMqttBroker();  // Connects to the configured mqtt broker
  setSubscriptions();     // Subscribe defined topics
}

void loop() {
  checkConnections(); // We check the connection every time
  pinta("p ");
}

/* Additional functions */
void setSubscriptions() {
  subscribe(TEST_TOPIC);   //Es subscriu al topic del broker que s'ha definit
}

void subscribe(char *newTopic) {
  const String topicStr = createTopic(newTopic);
  const char *topic = topicStr.c_str();
  mqttClient.subscribe(topic, QoS);  //Subscribes to messages published to the specified topic with QoS
  Serial.println("Client MQTT subscribed to topic: " + topicStr + " (QoS:" + String(QoS) + ")");
}

void callback(char *topic, byte *payload, unsigned int length) {   //procesa els missatges rebuts
  // Register all subscription topics
  static const String testTopicStr = createTopic(TEST_TOPIC);
  
  String msg = unwrapMessage(payload, length);
  Serial.println(" => " + String(topic) + ": " + msg);

  // Ordre que ha de fer depenent del topic
  if (String(topic) == testTopicStr) {
    pinta(msg);
  } else {
    Serial.println("[WARN] - '" + String(topic) +"' topic was correctly subscribed but not defined in the callback function");
  }
}

String unwrapMessage(byte *message, unsigned int length) {   //fa el missatge llegible
  String msg;
  for (int i = 0; i < length; i++) { // Unwraps the string message
    msg += (char)message[i];
  }
  return msg;
}

String createTopic(char *topic) {
  String topicStr = topic;
  return topicStr;
}

//Funció per conectar-se al WiFi i printejar-ho
void connectToWiFiNetwork() {
  Serial.print("Connecting with Wi-Fi: " + String(WIFI_SSID)); // Print the network which you want to connect
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);  //Inicialitza la configuració de red del la bibliotecta Wifi i proporciona l'estat actual
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".."); // Connecting effect
  }
  Serial.print("..connected!  (ip: "); // After being connected to a network,
                                       // our ESP32 should have a IP
  Serial.print(WiFi.localIP());
  Serial.println(")");
  String macAddressStr = WiFi.macAddress().c_str();
  strcpy(macAddress, macAddressStr.c_str());
}

//Funció per conectar-se al MQTT broker i printejar-ho
void connectToMqttBroker() {
  Serial.print("Connecting with MQTT Broker: " + String(MQTT_BROKER_IP));    // Print the broker which you want to connect
  mqttClient.connect(macAddress, MQTT_USER, MQTT_PASSWORD); //Connects the client. In our case we use the macAdress as the clientID
  while (!mqttClient.connected()) {   //Comproba si el client s'ha conectat al servidor (retorna false/true)
    delay(500);
    Serial.print("..");             // Connecting effect
    mqttClient.connect(macAddress); // Using unique mac address from ESP32
  }
  Serial.println("..connected! (ClientID: " + String(macAddress) + ")");
}

void checkConnections() {
  if (mqttClient.connected()) {
    mqttClient.loop();      //S'ha de cridar a la funció regularment per permetre que el client procesi els missatges entrants i mantingui la seva conexió amb el servidor.
  } else { // Try to reconnect (es reconnecta altra vegada al broker, al Wifi i es subscriu als topics corresponents
    Serial.println("Connection has been lost with MQTT Broker");
    if (WiFi.status() != WL_CONNECTED) { // Check wifi connection
      Serial.println("Connection has been lost with Wi-Fi");
      connectToWiFiNetwork(); // Reconnect Wifi
    }
    connectToMqttBroker(); // Reconnect Server MQTT Broker
    setSubscriptions();    // Subscribes to configured topics
  }
}

//Programa de LED

void pinta(String s) {
  //el string s s'ha de rebre del MQTT
  int l = s.length();         // Calcula llargada del string s  (l=2)
  /*
   El FOR serveix per tenir una pantalla negre (tots leds apagats) després de cada comanda
   'r' --> activa el RIGHT
   ' ' --> activa pantalla en negre (apagada)
   */
  for ( int n = 0; n< l; n++ )
    {
        long t = millis();  //Devuelve el número de milisegundos transcurrido desde el inicio del programa en Arduino hasta el momento actual
        char c = s[n];
        while ( millis()< t+ 400)
        SetChar(c);        //la funció dibuixa la comanda en pantalla
    }
}

bool GetBit( byte N, int pos)   //la funció permet desplaçar el que està en pantalla (anar-ho corrent)
   {                           // pos = 7 6 5 4 3 2 1 0
       int b = N >> pos ;      // movem la fila de bits cap a la dreta tantes posicions com el valor de pos
       b = b & 1 ;             // només fa la operació amb el últim bit b=001 --> b= 1&1 = 1 = TRUE 
       return b ;              // 1=TRUE i 0=FALSE
   }
   
byte Select( char c, byte row)  //row és un número entre 0 i 7
   {
       if ( c == ' ')          return(sp[row]);
       if ( c == 'r')          return(RIGHT[row]);
       if ( c == 'l')          return(LEFT[row]);
       if ( c == 's')          return(STOP[row]);
       if ( c == 'p')          return(RECTO[row]);
   }
    
void SetChar(char p)  
{  
  // iterate over the rows:
  for (int thisRow = 0; thisRow < 8; thisRow++) 
  {
    //take the row pin high:
    digitalWrite(row[thisRow], LOW);  //activem les files una per una
    byte F = Select( p, thisRow);     //F és la fila del caràcter p (que volem iluminar) corresponent 
    //iterate over the cols:
    for (int thisCol = 7; thisCol >=0; thisCol--) 
    {
      bool b = GetBit(F, thisCol);   //rertorna TRUE o FALSE
      if (b)  
      {
        digitalWrite(col[7-thisCol]  ,HIGH);   //activem la columna
        digitalWrite(col[7-thisCol]  ,LOW); //tornem a desactivar la columna
      }
      else 
      {
        digitalWrite(col[7-thisCol]  ,LOW);    // If it is 0, turn it off 
      }
    }
    // take the row pin low to turn off the whole row:
    digitalWrite(row[thisRow], HIGH); //tornem a desactivar la fila
  } 
}
