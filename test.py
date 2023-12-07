import subprocess
import multiprocessing
import time
import os

first_node = subprocess.Popen(['go', 'run', "./cmd/chord/.", "-a", "0.0.0.0", "-p", "2233"])
print("subprocess started")
nodes = []


num_of_nodes = 100
for i in range(100, 100+num_of_nodes):
    nodes.append(subprocess.Popen(['go', 'run', "./cmd/chord/.", "-ja", "0.0.0.0",  "-jp",  "2233", "-p", f"1{i}"]))
    time.sleep(0.1)

time.sleep(10)

print("done")
first_node.kill()
for node in nodes:
    node.kill()
    
exit(0)