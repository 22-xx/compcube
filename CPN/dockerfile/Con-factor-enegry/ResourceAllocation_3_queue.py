# -*- coding: utf-8 -*-
"""
Agorithm 1 of solving th optimal resource allocation problem (P4) in Sev.IV.B given offloading decisions in (P1)

Input: binary offloading mode

Output: the computation rate for all users

Created on May 2022

"""
import numpy as np
from scipy.optimize import linprog
from scipy.optimize import brent, fmin, minimize


# def Algo1_NUM(mode,h,w,Q,Y, V=20):
def Algo1_NUM(decion, price, dis, m, w, Q, H, Y, TTT, V, P_t,P_e):
    # print('QQQQ的数值',decion)   # N user, D datacenter,  N*D   10*2
    # print('QQQQ的数值',Q)   # N user, D datacenter, n time slot, data queue  H   N*D   10*2
    # print('YYYY的数值',Y)   # N user, D datacenter, n time slot, data queue  Q   N*D   10*2
    # print('price的数值',price)   # D datacenter, K VM,  D*M   2*3
    # print('dis的数值',dis)   # N users, D datacenter,   N*D   10*2
    # print('mmmm的数值',m)   # 1*60  1 random decision, N*M decision, 2 datacenter
    # print('wwww的数值',w)   # 10  workload, N users
    # P_t = 0.001  # 时间因素的权重因子
    # P_e = 0.01  # 能耗因素的权重因子
    phi = 100
    w_max = 2000

    N = 10  # 用户数量
    M = 4  # VM types
    D = 4  # datacenter number
    value_up_range = 0.82
    task_s_0 = np.zeros((N, D))
    VM = np.zeros((N, D))
    VM_SLA = np.zeros((N, D))
    price_p = np.zeros((D, N, M))

    # Demo：多变量边界约束优化问题(Scipy.optimize.minimize)
    # 定义目标函数
    def objf3(x):  # Rosenbrock 测试函数

        dis_delay = sum(sum(decion * Y * P_t * dis))
        VM_1 = m.reshape(D, N, M) * (x.reshape(D, N, M))

        for n in range(N):
            for d in range(D):
                # 计算数量*计算频率
                VM[n, d] = sum(VM_1[d, n, :])
                # 计算能耗
                VM_SLA[n, d] = sum(m.reshape(D, N, M)[d, n, :] * (x.reshape(D, N, M)[d, n, :] ** 3))
                # 扩展价格维数
                price_p[d, n, :] = price[d, :]
        for i in range(N):
            for d in range(D):
                task_s_0[i, d] = decion[i, d] * w[i]

        # 处理的任务
        rent = sum(sum(sum(m.reshape(D, N, M) * price_p)))
        # 来的任务
        task_s = sum(sum(Q * task_s_0))
        # 处理的任务
        task_l = sum(sum(Q * VM) / phi)
        SLA_T = sum(sum(H * P_t * (np.maximum(Q + (TTT - 1) * w_max / (np.ones((N, D))) - VM / phi, 0))))
        SLA_E = sum(sum(Y * P_e * VM_SLA))
        fx = dis_delay + task_s - task_l + SLA_T + SLA_E + V*rent

        return fx

    # 定义边界约束（优化变量的上下限）
    # b0 = (0.0, None)  # 0.0 <= x[0] <= Inf
    b = (0.6, 2)  # 0.0 <= x[1] <= 10.0
    bnds = tuple(b for x in range(0, D * N * M))
    xIni = np.ones((D * N * M)) * value_up_range
    cons = (
        {'type': 'ineq', 'fun': lambda x: 10 - x.reshape(D, N, M)[0, :, 0].sum(axis=0)},  # x>=e，即 x > 0
        {'type': 'ineq', 'fun': lambda x: 12 - x.reshape(D, N, M)[0, :, 1].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 13 - x.reshape(D, N, M)[0, :, 2].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 14 - x.reshape(D, N, M)[0, :, 3].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 12 - x.reshape(D, N, M)[1, :, 0].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 13 - x.reshape(D, N, M)[1, :, 1].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 12 - x.reshape(D, N, M)[1, :, 2].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 8 - x.reshape(D, N, M)[1, :, 3].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 10 - x.reshape(D, N, M)[2, :, 0].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 18 - x.reshape(D, N, M)[2, :, 1].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 12 - x.reshape(D, N, M)[2, :, 2].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 14 - x.reshape(D, N, M)[2, :, 3].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 8 - x.reshape(D, N, M)[3, :, 0].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 10 - x.reshape(D, N, M)[3, :, 1].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 12 - x.reshape(D, N, M)[3, :, 2].sum(axis=0)},
        {'type': 'ineq', 'fun': lambda x: 15 - x.reshape(D, N, M)[3, :, 3].sum(axis=0)},)

    # 优化计算
    xIni = np.ones((D * N * M)) * value_up_range
    resRosen = minimize(objf3, xIni, method='SLSQP', constraints=cons, bounds=bnds)
    xOpt = resRosen.x
    f_val = resRosen.fun

    # print('xOptxOpt',xOpt)

    # 最优的队列值
    VM_1 = m.reshape(D, N, M) * (xOpt.reshape(D, N, M))
    for n in range(N):
        for d in range(D):
            # 计算数量*计算频率
            VM[n, d] = sum(VM_1[d, n, :])
            # 计算能耗
            VM_SLA[n, d] = sum(m.reshape(D, N, M)[d, n, :] * (xOpt.reshape(D, N, M)[d, n, :] ** 3))
    # 最优值，data[i_idx:TTT-1,:,:],energy[i_idx:TTT-1,:,:], freque
    # 到了一定时间VM突然变大
    return f_val, VM, VM_SLA, xOpt.reshape(D, N, M)


