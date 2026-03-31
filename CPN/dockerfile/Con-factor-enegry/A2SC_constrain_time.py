#  #################################################################
#
#  This file contains the main code of Multi-resource-allocation.
#
# version 1.0 -- July 2022. Written by Xiaoxu Ren (xiaoxuren@tju.edu.cn)
#  #################################################################


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

myparams = {
    'axes.labelsize': '30',
    'xtick.labelsize': '25',
    'ytick.labelsize': '25',

    'lines.linewidth': 3.5,
    'legend.fontsize': '25',
    'font.family': 'Times New Roman'
}
pylab.rcParams.update(myparams)
plt.style.use('seaborn-deep')



# generate racian fading channel with power h and Line of sight ratio factor
# replace it with your own channel generations when necessary


def A2SC(n,SLA_time,P_t):
    '''
        LyDROO algorithm composed of four steps:
            1) 'Actor module'
            2) 'Critic module'
            3) 'Policy update module'
            4) ‘Queueing module’ of
    '''

    N =10                     # 用户数量
    K = N                   # initialize K = N
    M=4                  # VM types
    D=4                 # datacenter number
    J=20               # 总共的服务器节点
    Num_pro=10*N        # 每个区域的算力提供者
    VM_p=3              # VM parameter: CPU，memory, bandwidth
    Memory = 1024          # capacity of memory structure
    beta=0.1              # cost factor
    Reuest_time = 10
    # SLA_time = 0.4 # time comsumption threshold in J per time slot
    SLA_energy = 0.08  # energy comsumption threshold in J per time slot
    # P_t = 0.001  # 时间因素的权重因子
    P_e = 0.01  # 能耗因素的权重因子
    TTT=100
    V = 20
    # 用户的工作负载

    mem = MemoryDNN(net = [N*3, 256, 128, N],
                    learning_rate = 0.01,
                    training_interval=20,
                    batch_size=128,
                    memory_size=Memory
                    )

    start_time=time.time()
    mode_his = [] # store the offloading mode
    k_idx_his = [] # store the index of optimal provisioning actor
    price=np.zeros((n,D,M))


    freque=np.zeros((n,D,N,M))
    Q = np.zeros((n, N,D)) # N user, D datacenter, n time slot, data queue  H
    Y = np.zeros((n, N,D)) # N user, D datacenter, n time slot, SLA time queue   Q
    H = np.zeros((n, N, D))  # N user, D datacenter, n time slot, energy SLA queue   Q
    data = np.zeros((n, N,D)) # N user, D datacenter, n time slot, request decision
    delay=np.zeros((n, N,D))
    Obj = np.zeros(n) # objective values after solving problem (26)
    energy = np.zeros((n, N,D)) # energy consumption
    decision_old=np.zeros((n, N,D))
    decision_new=np.zeros((n, N,D))
    h=np.zeros((M,VM_p))   # VM parameter: utilizition, CPU, memery
    beta=np.zeros((M))   # VM service rate
    dis=np.zeros((N,D))
    workload = np.zeros((n,N))
    num = np.zeros((D,M))
    num_old=np.zeros((D*N*M))
    dis = dis_read()  # 生成用户和算力代理的距离


    for i_idx in range(n):
        print('i_idxi_idxi_idx',i_idx)


        # 六种虚拟机类型，CPU，存储，带宽, 考虑了五个不同的区域
        #100000,Q值可能后期为0
        h_new=VM_read()
        workload[i_idx,:]=20
        # workload[i_idx,:]=work_read(i_idx,N)**0.2

        beta=[random.uniform(0,1),random.uniform(0,1),random.uniform(0,1)]

        # 生成40个相同的虚拟机属性，后续需要更改 M*N个虚拟机属性  40*3种
        nn_input_old=[]


        # for d in range(D):
        #     for k in range(M):
        #         # price[d,k]=market_price(N,d,mode_his[i_idx-1])
        #         price[d,k]=random.randint(0,1)


        # 4) ‘Queueing module’ of LyDROO
        if i_idx > 0:
            # update queues
            # m_list[i_idx-1,k_idx_his[-1]] # 1-dim: user; 2-dim: datacenter; 3-dim: VM type, 4-dim: best

            for d in range(D):
                # 请求积压
                Q[i_idx,:,d] = np.maximum(Q[i_idx-1,:,d] - data[i_idx-1,:,d]+ decision_new[i_idx-1,:,d]*workload[i_idx-1,:],0)
                # print('data[i_idx-1,:,d]data[i_idx-1,:,d]',data[i_idx-1,:,d])
                # print('decision_new[i_idx-1,:,d]*workload[i_idx,:]/1024',decision_new[i_idx-1,:,d]*workload[i_idx-1,:])
                Q[i_idx,Q[i_idx,:,d]<0,d] =0                # assert Q is positive due to float error

                # SLA服务质量：时间
                delay[i_idx-1,:,d]=decision_new[i_idx-1,:,d]* dis[:,d] + np.maximum(Q[i_idx-1,:,d] - data[i_idx-1,:,d],0)
                H[i_idx, :, d] = np.maximum(H[i_idx - 1, :, d] + P_t * delay[i_idx - 1, :, d] - SLA_time,0)  # current energy queue
                H[i_idx, H[i_idx, :, d] < 0, d] = 0
                # print('P_t * delay[i_idx - 1, :, d]', P_t * delay[i_idx - 1, :, d])

                # 能耗队列m
                Y[i_idx, :, d] = np.maximum(Y[i_idx - 1, :, d] + P_e * energy[i_idx - 1, :, d] - SLA_energy,0)  # current energy queue
                # assert Y is positive due to float error
                Y[i_idx, Y[i_idx, :, d] < 0, d] = 0
                # print('P_e*energy[i_idx-1,:,d]', P_e*energy[i_idx-1,:,d])
                # 计算每一个用户在数据中心的权值  60S 进行一次任务调度

            for d in range(D):
                decision_old[i_idx , :, d] = workload[i_idx , :] * Q[i_idx , :, d] + H[i_idx , :, d] * P_t * dis[:, d]

            # 选择权值最小的数据中心，作为卸载策略
            for i in range(N):
                decision_new[i_idx, i, np.argmin(decision_old[i_idx, i, :])] = 1

        for d in range(D):
            data_old_queue=np.full((M,len(Q[i_idx,:,d])),Q[i_idx,:,d]/10000)   # (4,10)
            data_new_queue=np.vstack( ( data_old_queue)).transpose().flatten()  # (1,40)
            mig_old_queue=np.full((M,len(H[i_idx,:,d])),H[i_idx,:,d]/10000)   # (4,10)
            mig_new_queue=np.vstack( ( mig_old_queue)).transpose().flatten()  # (1,40)
            energy_old_queue = np.full((M, len(Y[i_idx, :, d])), Y[i_idx, :, d] / 10000)  # (4,10)
            energy_new_queue = np.vstack((mig_old_queue)).transpose().flatten()  # (1,40)
            nn_queue=np.vstack( ( data_new_queue,mig_new_queue,energy_new_queue)).transpose()   # (40,3)
            nn_input_old.append(nn_queue)
        mmmm = np.vstack((nn_input_old[0], nn_input_old[1],nn_input_old[2],nn_input_old[3]))
        nn_input =np.array( np.hstack((h_new, mmmm)))


        # 1) 'Actor module' of LyDROO
        # generate a batch of actions 整数内
        if (i_idx==0 or (i_idx+TTT)%TTT == 0):
            m_list = mem.decode(nn_input).transpose()   # 6*60
            r_list = [] # all results of candidate offloading modes
            v_list = [] # the objective values of candidate offloading modes
            # 生成VM类型在不同区域的价格 (i_idx+TTT)%TTT
            # D*M 每个区域所有用户的分配的不同虚拟机数量之和
            if i_idx==0:
                num=np.zeros((D,M))
            else:
                num=np.array(mode_his[int(i_idx/TTT-1)]).reshape(D,N,M).sum(axis=1)
                # num=num_old.reshape(D,N,M).sum(axis=1)
            # num=mode_price.sum(axis=1)
            # for d in range(D):
            #     for k in range(M):
            # price,u_bid_s,p_min_s,p_Off_s=market_price(i_idx,D,M,num,Num_pro)
            price[i_idx,:,:],u_bid_s,p_max_s,p_min_s,p_Off_s=market_price(i_idx,D,M,Reuest_time,num,Num_pro)

            for m in m_list:
              # 2) 'Critic module' of LyDROO
              # allocate resource for all generated offloading modes saved in m_list
              # 得到了最优值f, rate (N,D), VM_SLA (N,D)
              r_list.append(Algo1_NUM(decision_new[i_idx,:,:],price[i_idx,:,:],dis,m,workload[i_idx,:],Q[i_idx,:,:],H[i_idx,:,:],Y[i_idx,:,:],TTT, V,P_t,P_e))
              # 取出最优值f
              v_list.append(r_list[-1][0])

             # record the index of largest reward
            k_idx_his.append(np.argmin(v_list))
            # print('r_listr_list',np.shape(r_list))

            # 3) 'Policy update module' of LyDROO
            # encode the mode with largest reward
            mem.encode(nn_input, m_list[k_idx_his[-1]])
            # 200个策略，40*5 ，10个用户，4中虚拟机类型，5个区域
            mode_his.append(m_list[k_idx_his[-1]])
            # print('m_list[k_idx_his[-1]]m_list[k_idx_his[-1]]',np.shape(m_list[k_idx_his[-1]]))
            # print('type(r_list[k_idx_his[-1]])',r_list[k_idx_his[-1]])

            # 组合拍卖匹配
            # for d in range(D):
            #   for k in range(M):
            #     auction.model(N, J, provider_true_bid[d,m]*m_list[k_idx_his[-1]], capacity, user_true_bid)

            Obj[i_idx:i_idx+TTT-1], data[i_idx:(i_idx+TTT-1),:,:], energy[i_idx:(i_idx+TTT-1),:,:], freque[i_idx:(i_idx+TTT-1),:,:,:] = r_list[k_idx_his[-1]]
        # print('data[i_idx:TTT-1,:,:]',data[i_idx:TTT-1,:,:])

    Q = Q[:, :, :].sum(axis=1).sum(axis=1).sum(axis=0) / (n*N*D)
    H = H[:, :, :].sum(axis=1).sum(axis=1).sum(axis=0) / (n*N*D)
    Y = Y[:, :, :].sum(axis=1).sum(axis=1).sum(axis=0) / (n*N*D)

    return price,Q,H,Y




