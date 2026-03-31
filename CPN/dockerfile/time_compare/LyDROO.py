import scipy.io as sio                     # import scipy.io for .mat file I/
import numpy as np                         # import numpy

# for tensorflow2
from memoryTF2conv_3_queue import MemoryDNN
# from optimization import bisection
from ResourceAllocation_3_queue import Algo1_NUM
from game import VM_read,dis_read,market_price,work_read

import math

import time
import random
import matplotlib.pylab as pylab
import matplotlib.pyplot as plt
from LyDROO_TTT_queue_compare import MSL_DRL

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
obj = np.zeros((n))
cost = np.zeros((n))
# 设置不同的间隔
T_1=20
T_2=40
T_3=60
T_4=80
price_1,obj_1,Q_1,H_1,Y_1,cost_1=MSL_DRL(n,T_1)
price_2,obj_2,Q_2,H_2,Y_2,cost_2=MSL_DRL(n,T_2)
price_3,obj_3,Q_3,H_3,Y_3,cost_3=MSL_DRL(n,T_3)
price_4,obj_4,Q_4,H_4,Y_4,cost_4=MSL_DRL(n,T_4)




# fig = plt.figure(figsize=(30, 10))
#
# ax = plt.subplot(131)
# plt.plot(np.arange(0, n ), Q_1[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'y', linestyle = 'dashed', label='T=2')
# # plt.plot(np.arange(0, i_idx + 1), Q_near[:, :, :].sum(axis=1).sum(axis=1) / (N*D), label='Without dynamic price')
# plt.plot(np.arange(0, n ), Q_2[:, :, :].sum(axis=1).sum(axis=1) / (N * D),color = 'c', linestyle = 'dashed', label='T=20')
# plt.plot(np.arange(0, n ), Q_3[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'm', linestyle = 'dashed',label='T=50')
# plt.legend()
# plt.xlabel('Timeslot')
# plt.ylabel('Average Data Queue')
# # plt.savefig(r'Data_Queue.pdf', dpi=600)
#
# ax = plt.subplot(132)
# plt.plot(np.arange(0, n ), H_1[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'y', linestyle = 'dashed',)
# # plt.plot(np.arange(0, i_idx + 1), H_near[:, :, :].sum(axis=1).sum(axis=1) / (N * D), label='Without dynamic price')
# plt.plot(np.arange(0, n ), H_2[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'c', linestyle = 'dashed',)
# plt.plot(np.arange(0, n ), H_3[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'm', linestyle = 'dashed',)
# plt.xlabel('Timeslot')
# plt.ylabel('Average Time Queue')
#
#
# ax = plt.subplot(133)
# plt.plot(np.arange(0, n), Y_1[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'y', linestyle = 'dashed',)
# # plt.plot(np.arange(0, i_idx + 1), Y_near[:, :, :].sum(axis=1).sum(axis=1) / (N * D), label='Without dynamic price')
# plt.plot(np.arange(0,n), Y_2[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'c', linestyle = 'dashed',)
# plt.plot(np.arange(0,n), Y_3[:, :, :].sum(axis=1).sum(axis=1) / (N * D), color = 'm', linestyle = 'dashed',)
# plt.xlabel('Timeslot')
# plt.ylabel('Average Energy Queue')
#
# plt.savefig(r'E:\Kiki\Lyp+AC\Lyp+AC\figure\SLA_Queue_TT.pdf', dpi=600)

sio.savemat(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/node 38/time_TTT_data/0857/result_price_%d.mat' % T_1,{'price': price_1, 'objective': obj_1, 'data_queue': Q_1, 'time_queue': H_1, 'energy_queue': Y_1,'cost': cost_1})
sio.savemat(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/node 38/time_TTT_data/0857/result_price_%d.mat' % T_2,{'price': price_2,'objective': obj_2, 'data_queue': Q_2, 'time_queue': H_2, 'energy_queue': Y_2,'cost': cost_2})
sio.savemat(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/node 38/time_TTT_data/0857/result_price_%d.mat' % T_3,{'price': price_3,'objective': obj_3, 'data_queue': Q_3, 'time_queue': H_3, 'energy_queue': Y_3,'cost': cost_3})
sio.savemat(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/node 38/time_TTT_data/0857/result_price_%d.mat' % T_4,{'price': price_4,'objective': obj_4, 'data_queue': Q_4, 'time_queue': H_4, 'energy_queue': Y_4,'cost': cost_4})
