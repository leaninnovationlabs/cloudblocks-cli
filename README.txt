--------------------
Cloudblocks CLI
--------------------

Use:

cloudblocks execute --file=<file path>

Description: creates workload or updates existing workload and runs terraform apply

cloudblocks list

Description: Prints a list of active cloudblocks 

cloudblocks list-workloads

Description: Prints a list of active workloads 

cloudblocks

Description: Prints a help document with information about the cloudblocks cli

--------------------
Testing
--------------------
From root, run "go test ./..." to run all unit tests

--------------------
Install
--------------------
Dependencies: AWSCLI, Terraform, and Go 1.22+ are required to build and run Cloudblocks

To build from source, run "go build -o cloudblocks" from root 

--------------------
To do
--------------------
1. Create a script that moves the binary to /usr/local/bin or wherever and updates the $PATH environment variable so that
   cloudblocks_cli can be accessed from the terminal 
2. Add a feature that allows cloudblocks to read the variables.tf from a git repo where Terraform modules are located and parse that 
   into a JSON object so that the UI can display this to the user as configuration options
3. Find a place to store the work/ and  modules/ directories in the Unix filesystem  

4. format JSON when we write it 
5. see if we can convert JSON to YAML
6. see if we can add a modules refresh
7. make cloudblocks.json simple liek the package.json 
8. try to add a config file for each module. This will allow us to simplify cloudblocks.json. 
9. Think about a registry sort of system for modules
10. add default values for name, uuid, run_id 
   - edge case
11. bucket + region for s3 backend should be environment specific

12. see about making key for each workload in s3 bucket backend

13. make sure modules/ and work/ folders are generated w/ cloudblocks init if they don't already exist. Check main.go
14. make script to generate config for UI from variables.tf 
    - make it able to modify the module.json for the specific module. params

15. add dry-run that shows terraform plan. 

16. environments are important. add .env support 


- Simplify the modules structure 
        {
            "modules":
                {
                    "ec2": "0.0.1",
                    "s3": "0.0.13"
                }
        }
        each module will show a module.json file to define various module information
        {
            "type": "terraform",
            "version": "0.12.0",
            "description": "This is a simple terraform module that creates a VPC in AWS."
        }

- Provide the list of environments managed in config.json
    "envs":{
        "dev":{
            "bucket":"govpdfsandy",
            "region":"us-east-1"
        },
        "prod":{
            "bucket":"govpdfsandy",
            "region":"us-east-1"
        }
    }

- Add another working module (using cmdrunner)

- Utilize the env to deploy to corresponding environment
- Generate a unique runid if not passed and make sure it also works without any uuid passed
- Use Jinja templates
- See if we can remove the source from template.hcl files and derive based on the folder structure
- Dry run option for the cli
- Support file type parameters


---
UI work

1. Workload details screen
   - Workload detials
   - Each run details
   - Logs related to each run
2. Pull that from the cloud blocks table
3. Create the inputs based on the modules config provided in the table
4. Save workload params into the workload table to show it later on the UI
5. Add create by and updated by to all the tables
