# Possible future plans

This file contains ideas and concepts that have not been implemented yet.

## Kubernetes integration

In order for full kubernetes integration the system needs a backend that collects the information for us.

### HTTP API

The current idea is to create an HTTP API which will consume the results of a fuzzing test for multiple projects, aggregate them and show them to the end user.

Possible integrations:
- Slack notifications on fuzz crashes
- programatically pulling the github repositories and on change to master automactially fuzz (CI Integration)
- basic auth

### HTTP API Frontend

Since the results in an API would not be sufficient, a UI integration should be considered.
In order to not get to fancy here I suggest a regular gin/sql app that has useraccounts and groups as well es finegrained acces to every repo.

#### Users

Users should have a simple boostrap integration for account creation, pw reset as well as access control. Its also important to add jwt tokens in order for our ci system to be able to push to private projects

#### Fuzzing Overview

Index:

The index should deliver a list of repositories that are being monitored and their status. For instance:

- XXyCorp/vault:
  Fuzz Packages: 3 Fuzz status: No Panic

- XXyCorp/othervault:
  Fuzz Packages: 14 Fuzz status: Crashers found in 3 packages

DetailView:

The detailview shows the details of the last fuzz runs summarized for each packges:

- Helper/XOR
 Fuzz Functions: 7 Fuzz-Crashes in 7 Functions
  - Crasher 1
  - Crasher 2
  - Crasher 3

CrashersView:

Crasherview shows the details of the crash and the corresponding code. Syntax highlighting will need to be integrated for it to be usefull.
Underneath the crash report there should be a button "open jira ticket"

NotifyMe:

The notifyme page defines how you want to be notified of a crasher that has been found. I believe slack and email are good integrations. Basically by clicking a star button, the user subscribes to be notified over email.

NotifyAdmin:

Admins should also be able to select who automatically gets notified when a crasher occurs, basically autosubscribe.

## Simple s3 integration

Instead of being upload to a http endpoint, the data of a fuzz runn could also be stored on s3 in a specific folder and then later on be manually reviewed.
The question here is, who will watch the bucket for new crashers? Do we email from the tool directly? then we need to integrate smtp and configurations for that.

## Multiple sources in one config

Currently the tool relies on one config per source because of the possible container integration. The idea behind is that you mount a container with a config and it does its work.
The program could easily be adapted for multi source usage.