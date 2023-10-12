# goHive5 Overview

## Auth
### Login
| Description | gohive5
|:---|:---|
| Base hive object | CreateLogin()


## Case Management

### Alert
| Description | gohive5  |
|:---|:---|
| Create alert | CreateAlert() |
| Update alert | UpdateAlert() |
| Get alert | GetAlert() |
| Delete alert | DeleteAlert() |
| Find alerts timed by field | FindAlertsByFieldTimed() |
| Get alerts timed - gets alerts which were updated since a specific date | GetAlertsTimed() |
| Find alerts by custom field | FindAlertsByCustomField() |
| Merge alert to case | MergeAlert() |

### Case
| Description | gohive5  |
|:---|:---|
| Create case | CreateCase()  |
| Update case | UpdateCase() |
| Get case | GetCase() |
| Get cases timed | GetCasesTimed() |
| Get alerts associated with case | GetCaseAlerts()|
| Delete case | DeleteCase()  |
| Create case from alert | CreateCaseFromAlert() |
| Find case | FindCase() | 
| Find case by custom field | FindCaseByCustomField() | 
| Get Case status options (New/InProgress etc.) | GetCaseStatusOptions()|

## Comments

### Alert
| Description | gohive5  |
|:---|:---|
| Add alert comment | AddAlertComment() |
| Get alert comments | GetAlertComments() |

### Case
| Description | gohive5  |
|:---|:---|
| Add case comment | AddCaseComment() |
| Get case comment | GetCaseComments() |
| Get case comments since a specific date | GetCaseCommentsTimed() |


## Observables

### General
| Description | gohive5  |
|:---|:---|
| Get all observable types available | GetObservableTypes()|
| Get observable | GetObservable() |
| Update observable | UpdateObservable() |
| Delete observable | DeleteObservable() |


### Alert
| Description | gohive5  |
|:---|:---|
| Add observable | AddAlertObservable() |
| Get observables (all) | GetAlertObservables() |
| Get observable (single) | GetAlertObservable() |

### Case
| Description | gohive5  |
|:---|:---|
| Add observable | AddCaseObservable() |
| Get observables | GetCaseObservables() |
| Add file as an observable | AddCaseObservableFile() |
| Get observables filtered (dataType + value) | GetCaseObservablesFiltered() |

## Tasks

### General
| Description | gohive5  |
|:---|:---|
| Get task | GetTask() |
| Update task | UpdateTask() |
| Delete task | DeleteTask() |
| Add task log (Message underneath task) | CreateTaskLog() |
| Get task logs | GetTaskLogs() |

### Case
| Description | gohive5  |
|:---|:---|
| Add task | AddTaskToCase() |
| Get all tasks of a case | GetCaseTasks()|

## Template

### General
| Description | gohive5  |
|:---|:---|
| Get case template | GetCaseTemplate() |
| Delete case template | DeleteCaseTemplate() |
| Update case template | UpdateCaseTemplate() |

## Timeline

### General
| Description | gohive5  |
|:---|:---|
| Create customEvent | CreateTimelineEvent() |
| Get timeline | GetTimeline() |
| Get timeline event | GetTimelineEvent() |
| Delete timeline event | DeleteTimelineEvent() |
| Update timeline event | UpdateTimelineEvent()|

## TTP

### Alert
| Description | gohive5  |
|:---|:---|
| Add procedure to alert| AddAlertProcedure() |

### Case
| Description | gohive5  |
|:---|:---|
| Add procedure to case| AddCaseProcedure() |

### General
| Description | gohive5  |
|:---|:---|
| Get visible users | GetVisibleUsers() |
