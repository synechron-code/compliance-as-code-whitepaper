# Object Storage Compliance as code

This folder contains Cucumber tests which assert whether the controls around Object Storage are compliant with best practice.

Our examples use different techniques, as described in our whitepaper

## Controls Implemented

* [Encryption in flight](./encryption_in_flight/)
* [Encryption at rest](./encryption_at_rest/)
* [Restrict network access to known set of IP addresses](./access_whitelisting/)
 
## Techniques

The different Cloud Services Providers take materially different approaches towards achieving the same compliance effect with Object Storage. Below is a table stating the testable level of controls for each CSP.

We demonstrate various techniques - from the rudimentary to the more complex - to set-up the tests and attest to their efficacy

| Control Description | AWS | Azure|
|---|---|---|
|Encryption in Flight | Detective & Corrective | Preventative |
|Encryption at Rest | Self-Healing | *Not easily testable* |
|Restrict Network Access | Config Validation | Preventative |

For more detailed implementation information please see the respective README files.

## Future Developments

We also plan to build additional examples, to demonstrate how the ecosystem of tooling to support compliance activity in the cloud can be integrated with a common set of Behaviour Driven specifications and tests:

* Google Cloud Platform
* Hashicorp Sentinel
* Palo Alto Prisma Cloud