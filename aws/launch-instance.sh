 aws ec2 create-key-pair --key-name MyKeyPair --query 'KeyMaterial' --output text > ./MyKeyPair.pem
 chmod 400 MyKeyPair.pem
 aws ec2 create-security-group --group-name my-sg --description "My security group"
 aws ec2 run-instances --image-id ami-e90dc68a --count 1 --instance-type t2.micro --key-name MyKeyPair --security-groups my-sg

