# Cloudblocks CLI Tool

Cloudblocks is a command-line interface (CLI) tool for managing and interacting with cloud resource workloads. It provides a set of commands to initialize the environment, manage modules, execute workloads, and delete workloads. A workload is a single execution of a terraform module or script which creates cloud resources.

## Installation
Requirements: Linux or MacOS and Make

To install the Cloudblocks CLI tool, follow these steps:

From the root of the cloudblocks-cli directory, run the following commands:
1. make build
2. make install

## Usage

### Initializing the Cloudblocks Environment

Before using the Cloudblocks CLI tool, you need to initialize the cloudblocks environment. Run the following command:

```
cloudblocks init [--workdir=<work_directory>] [--modulesdir=<modules_directory>] [--bucket=<bucket_name>] [--region=<aws_region>]
```

This command creates the necessary directories and configuration files for cloudblocks. You can specify custom paths for the work directory and modules directory using the `--workdir` and `--modulesdir` flags, respectively. Additionally, you can set the S3 bucket name and AWS region using the `--bucket` and `--region` flags.

### Managing Cloudblock Modules

Cloudblocks uses modules to define the infrastructure components. You can manage cloudblock modules using the following commands:

#### Adding a Module

To add a new cloudblock module, use the following command:

```
cloudblocks modules add --name=<module_name> --version=<module_version>
```

This command adds a new module entry to the `modules.json` file with the specified name and version.

#### Updating a Module

To update an existing cloudblock module, use the following command:

```
cloudblocks modules update --name=<module_name> --version=<new_version>
```

This command updates the version of the specified module in the `modules.json` file.

#### Deleting a Module

To delete a cloudblock module, use the following command:

```
cloudblocks modules delete --name=<module_name>
```

This command removes the specified module entry from the `modules.json` file.

### Adding Environments

To add an environment to the config, use the following command:

``` 
cloudblocks env add --region --name --bucket
```

### Listing Environments

To list available environments and their configuration, use the following command: 

``` 
cloudblocks env list 
```

### Updating Environments 
To update an existing environment in the config, use the following command:

```
cloudblocks env update --region --name --bucket
```

### Deleting Environments
To delete an environment from the config, use the following command: 

```
cloudblocks env delete --region --name --bucket
```

#### Listing Modules

To list the available cloudblock modules, use the following command:

```
cloudblocks list
```

This command reads the content of the `modules.json` file and lists the available cloudblock modules in JSON format.

### Executing Workloads

To execute a workload, use the following command:

```
cloudblocks execute [<workload_json>] [--file=<workload_file>]
```

This command executes a workload by processing the user configuration into a `main.tf` file, creates a directory for the workload, and runs `terraform init` and `terraform apply`. You can provide the workload JSON either as a command-line argument or using the `--file` flag to specify the path to a JSON file containing the workload configuration.

### Dry-Run Workloads

To perform a dry-run of a workload, use the following command:

```
cloudblocks dry-run [<workload_json>] [--file=<workload_file>]
```

This command performs a dry-run of a workload by processing the user configuration into a `main.tf` file, creates a directory for the workload if one doesn't already exist, and runs `terraform init` and `terraform plan`. You can provide the workload JSON either as a command-line argument or using the `--file` flag to specify the path to a JSON file containing the workload configuration.

### Deleting Workloads

To delete a workload, use the following command:

```
cloudblocks delete [<workload_json>] [--file=<workload_file>]
```

This command runs `terraform destroy` to delete the resources associated with the workload and then deletes the workload directory and its contents. Similar to the execute command, you can provide the workload JSON either as a command-line argument or using the `--file` flag.

### Listing Workloads

To list the available workloads, use the following command:

```
cloudblocks list-workloads
```

This command lists the workloads defined in the `workloads.json` file.

## Configuration

The Cloudblocks CLI tool uses configuration files to store settings and module information. The main configuration files are:

- `config.json`: Stores the cloudblocks environment configuration, including the paths to the work directory and modules directory, as well as the environment information such as S3 backend bucket and AWS region.
- `modules.json`: Contains the list of available cloudblock modules, including their names and versions.
- `workloads.json`: Stores the list of executed workloads and their associated metadata.


## Building Go Binary
You can use the following command to build the go binary:
`go build -o tests/cloudblocks`


## Troubleshooting

If you encounter any issues while using the Cloudblocks CLI tool, consider the following:

- Ensure that you have initialized the cloudblocks environment using the `cloudblocks init` command before running other commands.
- Verify that you have the necessary permissions to access the specified S3 bucket and AWS region.
- Check the command usage and flags to ensure you are providing the correct arguments.
- If you encounter errors related to Terraform execution, review the Terraform logs for more information.

For further assistance, please refer to the Cloudblocks documentation or reach out to the Cloudblocks support team.
