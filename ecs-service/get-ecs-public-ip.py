import sys
import json
import subprocess
import time
import ipaddress


def get_ip_address_from_shell(lines, container_name, task_definition_name):

    for line in lines:
        is_matched_task = (str(line).rfind(task_definition_name) != -1)
        if is_matched_task:
            tokens = str(line).split()
            container_names = tokens[0].split("/")

            if len(container_names) > 1 and container_names[1].strip() == container_name:
                if tokens[1].strip() == 'RUNNING':
                    port_mappings = tokens[2].split(",")
                    # print(port_mappings)
                    if port_mappings:
                        ip_ports = port_mappings[0].split(":")
                        # print(ip_ports)
                        if ip_ports:
                            if not ipaddress.ip_address(str(ip_ports[0]).strip()).is_private:
                                return ip_ports[0]

    return None


def get_ip_address(cluster_config_name, container_name, task_definition_name):

    try:

        ip = None
        for index in range(30):

            p = subprocess.Popen(["ecs-cli", "ps", "--cluster-config", cluster_config_name],
                                 shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
            ip = get_ip_address_from_shell(p.stdout.readlines(), container_name, task_definition_name)

            if ip:
                break

            p.wait()

            if ip is None:
                time.sleep(3)

        data = {"public_ip": str(ip)}
        return data

    except Exception as e:
        sys.stderr.write("ERROR" + str(e))
        return None


def read_in():
    return {x.strip() for x in sys.stdin}


def main():

    # lines = read_in()

    file = sys.stdin.readline()
    # file = "{\"cluster_config_name\":\"next_telephony_2\"}"
    # print(file)
    data = json.loads(file)
    # print(data)
    cluster_config_name = data['cluster_config_name']
    container_name = data['container_name']
    task_definition_name = data['task_definition_name']

    output = get_ip_address(cluster_config_name, container_name, task_definition_name)
    if output:
        data = json.dumps(output)
        # sys.stdout.write('\'' + data + '\'')
        sys.stdout.write(data)
        return 0
    else:
        return 1


if __name__ == '__main__':
    sys.exit(main())
