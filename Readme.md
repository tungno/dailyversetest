## Om dette repositoriet 
Dette repositoriet er laget for å teste enhetstester for alle funksjonelle enheter i prosjektet. På grunn av begrenset tid har vi ikke rukket å implementere alle nødvendige tester, og derfor vil ikke alle enhetstester bestå foreløpig.


## Kjøring av enhetstester
For å kjøre alle enhetstestene fra rotmappen, kan man bruke følgende kommando:
```
 go test -cover ./...
```
Denne kommandoen kjører alle testfiler og viser dekningsgraden for testene.

## Kjøring av spesifikke tester
Hvis man ønsker å kjøre testene i en spesifikk mappe manuelt, kan man først navigere til den aktuelle mappen, for eksempel:
```
cd tests/handlers
```
Deretter kan man kjøre testene i den mappen direkte med 
```
go test <filenavn_test.go>
```
eller man kan gå til spesifikk fil og kjøre/teste på en funksjon. 