#!/bin/sh
#
# jmu, may 2021 
#
# create the workshop resources
#
# put your own values here
PASSWD=<robust password>
#
# this is the second region we are going to subscribe to
# we need regionname and region-key, see https://docs.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm
REG2=ap-sydney-1
REG2KEY=SYD
NAME=signals # dont change or find&replace in case you modify it!
# here the name of your home region
home=$(oci iam region-subscription list | jq -r '.data[0]."region-name"')
#
# create a compartment under tenancy
#
TENANCY=$(oci iam availability-domain list --all | jq -r '.data[0]."compartment-id"')
COMPARTMENT=$(oci iam compartment create -c $TENANCY --name $NAME --description $NAME | jq -r '.data.id')
sleep 5
echo "Created new compartement with ocid: "$COMPARTMENT
#
# subscribe to second region
#
echo "These are the regions already subscribed:"
oci iam region-subscription list | jq '.data[]."region-name"'
#
oci iam region-subscription create --tenancy-id $TENANCY --region-key $REG2KEY
#
# create dynamic group and iam policies needed for service connector, database resource principal and bucket replication
#
oci iam dynamic-group create --name $NAME --description $NAME --matching-rule "ALL {resource.type = 'autonomous-databases'}"
#
oci iam policy create --compartment-id $TENANCY --name signals --statements '["allow service objectstorage-eu-frankfurt-1 to manage object-family in compartment signals","allow any-user to {STREAM_READ, STREAM_CONSUME} in compartment signals","allow any-user to manage objects in compartment signals","allow dynamic-group signals to manage buckets in tenancy"]' --description signals
#
# object storage urls 
#
namespace=$(oci os ns get | jq -r .data)
echo "Grab this object storage urls for creating external tables later:"
echo "https://objectstorage.$home.oraclecloud.com/n/"$namespace"/b/signals1"
echo "https://objectstorage."$REG2".oraclecloud.com/n/"$namespace"/b/signals2"
#
# create databases and outputs sqldev urls
#
echo "Creating databases, grab the sqldeveloper url's for later:"
oci db autonomous-database create --display-name db1 --db-name db1 --cpu-core-count 1 --compartment-id $COMPARTMENT --admin-password $PASSWD --data-storage-size-in-tbs 1 --region $home --wait-for-state AVAILABLE | jq '.data."connection-urls"."sql-dev-web-url"'
#
oci db autonomous-database create --display-name db2 --db-name db2 --cpu-core-count 1 --compartment-id $COMPARTMENT --admin-password $PASSWD --data-storage-size-in-tbs 1 --region $REG2 --wait-for-state AVAILABLE | jq '.data."connection-urls"."sql-dev-web-url"'
#
# create streaming topic and output its ocid for later
#
echo "Creating streaming topic, grab the ocid of the streaming topic for later use in the microservice code:"
strocid=$(oci streaming admin stream create --name signals1 --partitions 1 --compartment-id $COMPARTMENT --region $home | jq -r '.data.id')
#
# generating json config for later create the service connector
echo '{"kind":"streaming","streamId":'$strocid',"cursor":{"kind":"LATEST"}}' > source
echo '{"kind": "objectStorage","bucketName": "signals1","objectNamePrefix": "signalsrepl"}' >target
#
# create os buckets
#
echo "Creating storage buckets in both regions"
oci os bucket create --name signals1  --compartment-id $COMPARTMENT --region $home 
oci os bucket create --name signals2  --compartment-id $COMPARTMENT --region $REG2 
#
# create service connector hub
#
echo "Creating service connector"
oci sch service-connector create --display-name signals --compartment-id $COMPARTMENT --region $home --source file://./source --target file://./target
#
# create bucket replication policy
#
echo "Creating object storage replication policy"
 oci os replication create-replication-policy  --bucket-name signals1 --name signalsrepl --destination-region $REG2 --destination-bucket signals2
 #



