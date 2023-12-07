import subprocess
import multiprocessing
import signal
import time
import os


test_output = open("test_output.txt", "w")
test_output.truncate(0)

first_node = subprocess.Popen(['go', 'run', './cmd/chord/.', '-a', '0.0.0.0', '-p', '2233'], shell=True, stderr=test_output, stdout=test_output, stdin=subprocess.PIPE)
time.sleep(2)

first_node.stdin.write("set a 123\n".encode())
first_node.stdin.write("get a\n".encode())

nodes = []
num_of_nodes = 5
for i in range(100, 100+num_of_nodes):
    nodes.append(subprocess.Popen(['go', 'run', "./cmd/chord/.", "-ja", "0.0.0.0",  "-jp",  "2233", "-a", "0.0.0.0", "-p", f"1{i}"], shell=True, stdin=subprocess.PIPE, stderr=test_output, stdout=test_output))
    time.sleep(0.1)

print("done creating nodes")
time.sleep(3)
nodes[0].stdin.write("set b 456\n".encode())

time.sleep(3)
nodes[-1].stdin.write("get b\n".encode())

time.sleep(1)
print("exiting")
test_output.close()
for node in nodes:
    node.send_signal(signal.CTRL_C_EVENT) 
first_node.send_signal(signal.CTRL_C_EVENT)

# time.sleep(5)
# os.system('type test_output.txt | findstr /v "Enter command: " > temp.txt && move /y temp.txt test_output.txt')
exit()

