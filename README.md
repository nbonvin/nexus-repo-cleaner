# nexus-repo-cleaner
Small program to clean nexus repositories

## Installation

```
go get github.com/nbonvin/nexus-repo-cleaner
```

## Examples 

Keep 10 most recent maven SNAPSHOT versions (default):
```
./nexus-repo-cleaner /opt/nexus/storage/snapshots/ 
```

Keep 20 most recent maven SNAPSHOT versions:
```
./nexus-repo-cleaner -keep=20 /opt/nexus/storage/snapshots/ 
```

Keep 20 most recent maven RELEASE versions:
```
./nexus-repo-cleaner -keep=20 /opt/nexus/storage/releases/ 
```
