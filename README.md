# Overview

This project implements a web-based UI for healthcare wristbands. It is implemented as a Golang server using a ReactJS client with a GraphQL data API.

# Unfinished Stuff & Bug

## Inconsistent Data Returned from Subscription
GQLPlayground is able to receive data quickly and smoothly whereas our React frontend struggles to keep up with the changes and sometime loses packets along the way of transporting data. This behaviour particularly shows in the nurse's main interface where the Alert subscription used for the main view is initally slower than other subscriptions such as: Alert subscription in the Vital Alerts section and the Summary subscription. Then, during the runtime of the system, it actually catches up and sometimes goes faster than the other Alert subscription

## No "Other section" in Vital Alerts
Vital Alerts doesn’t have "Other section" yet for under 16, pregnant and error bands. The bands which have errors are not implemented to be included in the High level alert either

## Data Not Coming Through When Using 2 Windows Concurrently
Neither running both windows on the same port (3000) or different ports (3000 and 8080) would create realtime communications. Data seemed to come through much slower when two windows are being used concurrently

# Showcase

## Folders
To play with the showcase, go to main.go file, find the finalFolder variable. There you can change the variable between Nurse and Admin, that way, the system will try and look for according folder with corresponding JSON files

## Vital Changes Priority
Based on descending alphabetical order (A, B, C, D, ... , Z)

# Project Structure

## Go
Project structure broadly follows https://github.com/golang-standards/project-layout with all code in "internal" rather than "pkg" as this is never meant to be included in another project.

run 'go mod vendor' after cloning project to get all the dependencies.

## React


# Code generation

Code generation is used in this project to create the GraphQL interfaces, both on the server and in the frontend React, from a shared schema.

## Go GraphQL generation
The Go code generation is integrated into the go build system:<br />
    Run either of the two command lines below to generate Golang code that handles GraphQL
    go run scripts/gqlgen.go

    Or

    go generate ./...

    N.B. you must have run 'go mod vendor' first.

It is implemented using github.com/99designs/gqlgen and configured using gqlgen.yml
The integration of gqlgen into go generate uses scripts/gqlgen.go.


## React GraphQL generation

# GNU Make
To install GNU MAKE https://chocolatey.org/packages/make install Chocolatey then run the following command in cmd:

```bash
choco install make
```

# Testing

## Go Unittests
To run all tests:
    go test ./...
To run an individual test, e.g.:
    go test -v -timeout 30s -run ^TestNewWristbandProxyFile$ cerberus-security-laboratories/des-wristband-ui/internal/core

## End-to-End tests


# Building and Running
This project utilise a Makefile to build and run the files.

## Run the project's frontend and backend concurrently on Golang

To re-build the ui folder, run the command:

```bash
make build (Requires GNU Make)
```

Moving onto the frontend code prep, you will need to make sure all Nodejs packages are installed by, from the root folder, running:

```bash
cd ui
```
```bash
npm install
```

Once thats done, run:

```bash
npm run build
```

This will build the static file for the Go server to be ran.

Finally, run the environment with the command:

```bash
make run (Requires GNU Make)
```

## Run the project's frontend and backend separately on Golang and React
(while in Development of the Front-End though, you may run the Go Server and another React App seperately to test GraphQL without having to npm run build to test every time)

### Client
Run the following command line to get the client up and running

```bash
npm start
```

### Server
go generate ./...

To run the server, either of these two command lines below should work

```bash
go run cmd/main.go
```

```bash
make run (Requires GNU Make)
```

Options:
    --ip and --port define the address and port of the webserver, defaults to localhost:8080

    --data-path and --data-prefix define the loaction of the wristband data files that need to be read in.
    e.g. for files wbData_001.json and wbData_002.json in tests/e2e/no_warnings use 
        'go run cmd/main.go --data-path tests/e2e/no_warnings --data-prefix wdData_'
    
    --tick defines the interval in milliseconds between reading data entries from the wristband files.
        For realistic behaviour use 10000, for fast simulation use 500 or 1000.

# Available Features on The GQLPlayground

## Mutations

### Adding New Wristband

```
mutation {
  addWristband(input: {
    tic: "2",
    name: "Doan Trong Vu",
    dateOfBirth: "17-08-1998"
    onOxygen: true,
    pregnant: false,
    child: true,
    department: "here",
  }) {
    id, 
    tic,
    active,
    name,
    dateOfBirth,
    onOxygen,
    pregnant,
    child,
    key,
    department,
  }
}
```

### Deactivating A Wristband

```
mutation {
  deactivateWristband(input: {
    id: 1
  }) {
    id,
    tic,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
  }
}
```

### Reactivating A Wristband

```
mutation {
  reactivateWristband(input: {
    id: 1
  }) {
    id,
    tic,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
  }
}
```

### Resetting

```
mutation {
  resetName (id: 1, value: "Elax Hipkons") {id}
}
```

```
mutation {
  resetDepartment (id: 1, value: "Over There") {id}
}
```

You can reset other wristband's input parameters separately like this, or you can choose to reset multiple fields by:
```
mutation {
  resetMultipleFields (id: 1, options: ["Name", "OnOxygen", "Child", "Department"], values: ["Vu Doooan", "false", "false", "There"])
  {
    id,
    tic,
    active,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
  }
}
```

**__id__ IS TO IDENTIFY WHICH BAND DATA SHOULD BE CHANGED**<br />
**__options__ IS TO IDENTIFY WHICH FIELDS ARE GONNG BE CHANGED (1: Name, 2: OnOxygen, 3: Pregnant, 4: Child, 5: Department)**<br />
**__values__ IS TO REPLACE THOSE OLD VALUES WITH THE NEW ONES, ACCORDINGLY TO THE OPTIONS**

### Reassigning to a New Wristband

```
mutation reassignNewWristband ($oldWristband: DeactivateWristbandInput!) {
  reassignNewWristband(oldWristband: $oldWristband) {
    id,
    tic,
    active,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
  }
}
```

## Queries

### Getting All Available Wristbands

```
query getWristbands {
  getWristbands {
		id,
    tic,
    active,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
    data {
      sensorData {
        respiration,
        sp02,
        pulse,
        temperature,
        bloodPressure
        motion,
        proximity
      },
      batteryLevel
    }
  }
}
```

### Getting A Wristband

```
query getWristband ($id: ID!){
	getWristband (id: $id){
		id,
    tic,
    active,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
    data {
      sensorData {
        respiration,
        sp02,
        motion,
        temperature
      },
      batteryLevel
    }
  }
}
```

### Getting A Wristband's Data

```
query getWristbandData ($id: ID!, $how: String!) {
  getWristbandData (id: $id, how: $how) {
    sensorData {
      respiration,
      sp02,
      pulse,
      temperature,
      bloodPressure
      motion,
      proximity
    },
    batteryLevel
  }
}
```

**__how: "first" or "latest"__ IS USED TO GET THE FIRST OR THE LATEST DATA IN THE TREND**

### Getting Multiple Wristband's Data

```
query getMultipleWristbandData ($id: ID!, $howMany: Int!, $start: Int!, $end: Int!) {
  getMultipleWristbandData (id: $id, howMany: $howMany, start: $start, end: $end) {
    sensorData {
      respiration,
      sp02,
      pulse,
      temperature,
      bloodPressure
      motion,
      proximity
    },
    batteryLevel
  }
}
```

**__howMany: 0, start: [Any], end: [Any]:__ Get All**<br />
**__howMany: > 0, start: [Any], end: [Any]:__ Get A Certain Amount of Data (Backwards in Time)**<br />
**__howMany: -1, start: [A Number], end: [A Number]:__ Switch The Nature of Querying to Getting Block of Wristband Data (Backwards in Time)**

## Subscriptions

### Updating Wristband

```
subscription updateWristbandAdded {
  updateWristbandAdded {
    id
    tic
    active
    name
    dateOfBirth
    onOxygen
    pregnant
    child
    key
    department
  }
}
```

### Updating Wristband Data

```
subscription updateWristbandData ($id: ID!) {
  updateWristbandData (id: $id) {
    sensorData {
      respiration,
      sp02,
      pulse,
      temperature,
      bloodPressure
      motion,
      proximity
    },
    news2 {
      sp02,
      pulse,
      temperature,
      bloodPressure
      motion,
      onOxygen,
      overall
    },
    batteryLevel
  }
}
```

### Updating Wristband Data Alert (For the Alert section)

```
subscription updateWristbandDataAlert($id: ID!) {
  updateWristbandDataAlert(id: $id) {
    sensorData {
      respiration,
      sp02,
      pulse,
      temperature,
      bloodPressure
      motion,
      proximity,
    },
    level,
    target,
    overall
  }
}
```

### Updating Wristband Data Summary

```
subscription updateSummary {
  updateSummary {
    high,
    medium,
    lowMedium,
    low,
    other
  }
}
```

### Updating Important Bands

```
subscription updateImportantBands {
  updateImportantBands {
    id,
    tic,
    active,
    name,
    onOxygen,
    pregnant,
    child,
    key,
    department,
    data {
      sensorData {
        respiration,
        sp02,
        motion,
        temperature
      },
      news2 {
        respiration,
        sp02,
        motion,
        temperature,
        motion,
        overall
      }
      batteryLevel
    }
  }
}
```
