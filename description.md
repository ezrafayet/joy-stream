Objectif du projet

Permettre à un joueur distant (au Japon) de jouer en multijoueur local sur une Nintendo Switch physique située chez toi, sans émulateur.
La Switch doit croire qu’une manette Bluetooth normale est connectée.

Le système utilise JoyControl pour émuler une manette Switch.

Architecture générale :

Manette du joueur distant
        ↓
client réseau
        ↓ UDP
serveur réseau (PC Linux)
        ↓
bridge JoyControl
        ↓ Bluetooth
Nintendo Switch

La Switch voit simplement un Pro Controller Bluetooth.

Architecture logicielle

Le projet contient 3 composants.

1. Remote Controller Client

Programme exécuté sur le PC du joueur distant.

Responsabilités :

lire une manette locale

envoyer l’état de la manette 60 fois par seconde

protocole UDP très léger

Bibliothèque recommandée :

pygame pour lire la manette

Envoi vers l’IP du serveur.

2. UDP Input Server (Go)

Programme principal exécuté sur ton PC Linux.

Responsabilités :

écouter sur UDP

décoder les paquets de contrôle

maintenir l’état courant de la manette

transmettre cet état au bridge JoyControl

Ce programme doit être écrit en Go.

3. JoyControl Bridge (Python)

Petit programme Python qui :

démarre un Pro Controller virtuel

attend l’appairage avec la Switch

reçoit les états de manette du serveur Go

met à jour l’état JoyControl

envoie les états Bluetooth à la Switch

Il utilise directement l’API interne de JoyControl.

Protocole réseau

Communication :

UDP
60 Hz
12 bytes par paquet

Format du paquet :

Byte 0-1   sequence number (uint16 big endian)
Byte 2-3   button bitmask (uint16)
Byte 4     LX
Byte 5     LY
Byte 6     RX
Byte 7     RY
Byte 8     dpad
Byte 9     misc flags
Byte 10-11 reserved
Mapping boutons
bit0  A
bit1  B
bit2  X
bit3  Y
bit4  L
bit5  R
bit6  ZL
bit7  ZR
bit8  PLUS
bit9  MINUS
bit10 LCLICK
bit11 RCLICK
bit12 HOME
bit13 CAPTURE
Axes

Valeurs :

0..255
128 = centre

Conversion :

float = (value - 128) / 127
Fréquence

Le client envoie :

60 packets / seconde

Le serveur met à jour l’état JoyControl à la même fréquence.

Fonctionnement côté Switch

Sur la Switch :

Settings
Controllers
Change Grip/Order

Le contrôleur JoyControl apparaît comme :

Pro Controller

La Switch pense que c’est une manette normale.

Structure du projet
remote-switch-play/
 ├── client/
 │    controller_client.py
 │
 ├── server/
 │    input_server.go
 │
 ├── bridge/
 │    joycontrol_bridge.py
 │
 └── protocol/
      packet.go
Flux complet
manette joueur
     ↓
client pygame
     ↓ UDP 60Hz
serveur Go
     ↓
bridge JoyControl
     ↓ Bluetooth
Nintendo Switch
Contraintes importantes

Le système doit :

utiliser UDP

tolérer la perte de paquets

utiliser des snapshots d’état (pas des événements)

ignorer les paquets avec numéro de séquence ancien

supporter latence ~200 ms

Comportement attendu

Quand le joueur appuie sur A :

client envoie paquet UDP

serveur Go met à jour l’état

bridge Python met à jour JoyControl

JoyControl envoie Bluetooth

la Switch reçoit l’input

Le jeu voit simplement une manette normale.

Objectif final

Pouvoir jouer :

toi (manette locale)
+
ton frère (manette distante)

dans un jeu multijoueur local Switch.

Simplifications autorisées

pas besoin de chiffrement

pas besoin de TCP

pas besoin de compression

pas besoin de gestion multi-joueurs