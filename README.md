# meshed



### Motivation

Sehr oft habe ich während meiner Tätigkeit Systeme gesehen, deren Entwickler zuerst über die Ablage der Daten und erst danach über die eigentlichen Probleme nachgedacht haben. Während der Lebenszeit dieser Systeme sind oft die ursprünglichen Ideen verloren gegangen und die Datenhaltung entwickelte sich immer mehr zum Hemmnis.
Diesen Zustand habe ich immer als sehr unbefriedigend empfunden.

Mit diesem "Mesh of Informations" gehe ich der Idee nach, mit einem toteinfachen Datenmodell beliebige Infomationen miteinander zu verknüpfen und ständig rekonfigurierbar zu halten.

### Idee

Die Grundidee ist es, ein Mesh bereitzustellen, in dem Informationen geregelt miteinander vernetzt und damit in fachliche Zusammenhänge gebracht werden können. 

Die Art und Weise der Verknüfung ist zur Lauf- und Lebenszeit des Systems einfach rekonfigurierbar, ohne das die Informationen selbst angepasst werden müssen.

Es entsteht ein Netzwerk an Informationen
- alle Informationen sind an Knoten gebunden
- von Knoten gehaltene Informationen sind typisiert
- Knoten können nach Regeln miteinander verknüpft werden
- Knoten enthalten Metainformationen über die Informationen selbst (History), die Erzeugung dieser Metadaten erfolgt vom Entwicker unabhängig

###Vorteile
- konkrete Implementationen von Programmen auf der Basis des Meshs unterscheiden sich nur über die enthaltenen Typen und deren Regeln, die Art und Weise der Verknüpfung erfolgt immer gleich
- unabhangig von der verwendeten Datenbank (was in dem vorliegenden Prototypen leider noch nicht vollständig gelungen ist)
- jeder einzelne implementierte Typ bestimmt selbst die Regeln, mit denen er mit anderen Typen vernetzt werden kann

###Was leistet der Prototyp

- Mesh und Nodes als Informationsträger sind beispielhaft implementiert
- beispielhaft wurden Informationen mit den Typen User, Image, Category implementiert
- Bilder können Usern zugeordnet werden, Bilder und User können mit Hilfe des Typs Category kategorisiert werden
- beispielhaftes Ändern der Informationen selbst über z.B. User, deren Passwort geändert wird
- beispielhaftes Veröffentlichen von Informationen über die im Mesh enthaltenen Typen und Abruf derer über eine API
- beispielhaftes Registrieren von Usern
- Verwendung von MongoDB als Persistenzschicht
- Verwendung von Mongo GridFS zur Ablage von Images

###Nachvollziehen der Funktionsweise

- zu jedem Typ und zur API wurden Tests implementiert

Translate to english -> tbd.
