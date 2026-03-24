import matplotlib.pyplot as plt
import matplotlib.pylab as pylab
from mpl_toolkits.mplot3d import Axes3D
import numpy as np
import scipy.io as sio
import pandas as pd
from matplotlib import cm

fig = plt.figure()

import seaborn as sns

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
# plt.style.use(['ggplot','seaborn-deep'])


fig4 = plt.figure()
ax4 = plt.axes(projection='3d')

data0 = sio.loadmat(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/Lyp+AC+2_3_queue/factor_compare/constrain/pic/result_constrain_SLA.mat')

print(data0.keys())
# p_t = np.arange(0.0001,0.01,0.001)
# p_e = np.arange(0.001,0.1,0.01)

SLA_time = np.arange(0.01,1,0.1)
P_t = np.arange(0.0001,0.01,0.001)
X, Y = np.meshgrid(SLA_time, P_t)

# print(data.keys())
matrix0_Q = data0['data_queue']
print(np.shape(matrix0_Q))
Z=matrix0_Q.reshape(10,10)
print(Z)

print(np.shape(Z))
print('Z Z Z ',Z )

surface = ax4.plot_surface(X,Y,Z, cmap=plt.cm.viridis, linewidth=0.2)
fig.colorbar(surface)
plt.savefig(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/Lyp+AC+2_3_queue/factor_compare/data_3D_SLA_energy.pdf', dpi=600)



fig2 = plt.figure()
ax2 = plt.axes(projection='3d')
# p_t = np.arange(0.0001,0.01,0.001)
# p_e = np.arange(0.001,0.1,0.01)
SLA_time = np.arange(0.01,1,0.1)
P_t = np.arange(0.0001,0.01,0.001)
X, Y = np.meshgrid(SLA_time, P_t)

# print(data.keys())
matrix0_H = data0['time_queue']
Z_H=matrix0_H.reshape(10,10)
print(Z_H)
surface = ax2.plot_surface(X,Y,Z_H, cmap=plt.cm.viridis,  antialiased=False)
fig.colorbar(surface)
plt.savefig(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/Lyp+AC+2_3_queue/factor_compare/time_3D_SLA_energy.pdf', dpi=600)


fig3 = plt.figure()
ax3 = plt.axes(projection='3d')
SLA_time = np.arange(0.01,1,0.1)
P_t = np.arange(0.0001,0.01,0.001)
X, Y = np.meshgrid(SLA_time, P_t)

# print(data.keys())
matrix0_Y = data0['energy_queue']
Z_Y=matrix0_Y.reshape(10,10)
print(Z_Y)


surface = ax3.plot_surface(X,Y,Z_Y,cmap=plt.cm.viridis,  antialiased=False)
fig.colorbar(surface)
plt.savefig(r'/Users/renxiaoxu/NutstoreCloudBridge/我的坚果云/要备份！！！！！/撰写论文/1115-Info/info2023/Lyp+AC+2_3_queue/factor_compare/energy_3D_SLA_energy.pdf', dpi=600)


# ax4.legend()
plt.show()

#作图
# ax = plt.subplot(131)
# cmap=cm.coolwarm,
# ax4.plot_surface(X,Y,Z,alpha=1,cmap='rainbow')     #生成表面， alpha 用于控制透明度
# ax4.contour(X,Y,Z,zdir='z', offset=-3,cmap="rainbow")  #生成z方向投影，投到x-y平面
# ax4.contour(X,Y,Z,zdir='x', offset=-6,cmap="rainbow")  #生成x方向投影，投到y-z平面
# ax4.contour(X,Y,Z,zdir='y', offset=6,cmap="rainbow")   #生成y方向投影，投到x-z平面
# ax4.contourf(X,Y,Z,zdir='y', offset=6,cmap="rainbow")   #生成y方向投影填充，投到x-z平面，contourf()函数

#设定显示范围
# ax4.set_xlabel('X')
# ax4.set_xlim(0, 0.005)  #拉开坐标轴范围显示投影
# ax4.set_ylabel('Y')
# ax4.set_ylim(0, 0.05)
# ax4.set_zlabel('Z')
# ax4.set_zlim(0, 3)

#
# matrix1_Q = data1['data_queue']