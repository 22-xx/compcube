import scipy.io as sio                     # import scipy.io for .mat file I/
import numpy as np                         # import numpy

# for tensorflow2
from memoryTF2conv_3_queue import MemoryDNN
# from optimization import bisection
from ResourceAllocation_3_queue import Algo1_NUM

import math

import time
import random
import matplotlib.pylab as pylab
import matplotlib.pyplot as plt
from A2SC_constrain_energy import A2SC
from matplotlib import pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import numpy as np

#定义坐标轴
fig4 = plt.figure()
ax4 = plt.axes(projection='3d')

myparams = {
    'axes.labelsize': '20',
    'xtick.labelsize': '20',
    'ytick.labelsize': '20',

    'lines.linewidth': 2,
    'legend.fontsize': '20',
    'font.family': 'Times New Roman',

}
pylab.rcParams.update(myparams)
plt.style.use('seaborn-muted')

N = 10  # 用户数量
n=500
K = N  # initialize K = N
M = 4  # VM types
D = 4  # datacenter number
J = 20  # 总共的服务器节点
Q = np.zeros((n, N, D))  # N user, D datacenter, n time slot, data queue  H
Y = np.zeros((n, N, D))  # N user, D datacenter, n time slot, SLA time queue   Q
H = np.zeros((n, N, D))  # N user, D datacenter, n time slot, energy SLA queue   Q
price = np.zeros((n, D, M))
# 设置不同的间隔
Q_1_total=[]
Q_2_total=[]
Q_3_total=[]
Q_4_total=[]

H_1_total=[]
H_2_total=[]
H_3_total=[]
H_4_total=[]

Y_1_total=[]
Y_2_total=[]
Y_3_total=[]
Y_4_total=[]

# p_t = np.arange(0.0001,0.01,0.001)
# p_e = np.arange(0.001,0.1,0.01)

SLA_energy = np.arange(0.01,1,0.05)
P_e = np.arange(0.001,1,0.1)

print(np.shape(SLA_energy))
print(np.shape(P_e))

# p_e_2=0.05
# p_e_3=0.1
# p_e_4=0.15
# p_e_5=0.2
#
X, Y = np.meshgrid(SLA_energy, P_e)

for i in range(len(SLA_energy)):
    for j in range(len(P_e)):
        print(i,j)
        price_1,Q_1,H_1,Y_1=A2SC(n,SLA_energy[i],P_e[j])
        Q_1_total.append(Q_1)
        H_1_total.append(H_1)
        Y_1_total.append(Y_1)
        # price_2,Q_2,H_2,Y_2=MSL_DRL(n,p_t[i],p_e_2)
        # Q_2_total.append(Q_1)
        # H_2_total.append(H_1)
        # Y_2_total.append(Y_1)
        # price_3,Q_3,H_3,Y_3=MSL_DRL(n,p_t[i],p_e_3)
        # Q_3_total.append(Q_1)
        # H_3_total.append(H_1)
        # Y_3_total.append(Y_1)
        # price_4,Q_4,H_4,Y_4=MSL_DRL(n,p_t[i],p_e_4)
        # Q_4_total.append(Q_1)
        # H_4_total.append(H_1)
        # Y_4_total.append(Y_1)



# price_2,Q_2,H_2,Y_2=MSL_DRL(n,p_t_1)
# price_3,Q_3,H_3,Y_3=MSL_DRL(n,p_t_1)


sio.savemat(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/Lyp+AC+2_3_queue/factor_compare/constrain/pic/result_con_energy_SLA.mat', {'price': price_1, 'data_queue': Q_1_total, 'time_queue': H_1_total, 'energy_queue': Y_1_total})
# sio.savemat(r'E:\Kiki\Lyp+AC\Lyp+AC\figure\result_factor_%d.mat' % p_e_2,{'price': price_2, 'data_queue': Q_2_total, 'time_queue': H_2_total, 'energy_queue': Y_2_total})
# sio.savemat(r'E:\Kiki\Lyp+AC\Lyp+AC\figure\result_factor_%d.mat' % p_e_3,{'price': price_3, 'data_queue': Q_3_total, 'time_queue': H_3_total, 'energy_queue': Y_3_total})
# sio.savemat(r'E:\Kiki\Lyp+AC\Lyp+AC\figure\result_factor_%d.mat' % p_e_4,{'price': price_4, 'data_queue': Q_4_total, 'time_queue': H_4_total, 'energy_queue': Y_4_total})
