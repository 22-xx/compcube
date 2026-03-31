#  #################################################################
#
#  This file contains the main code of LyDROO.
#
#  #################################################################
import random
import pandas as pd
import csv


import scipy.io as sio                     # import scipy.io for .mat file I/
import numpy as np                         # import numpy
import datetime as dt
import pandas as pd
from math import radians, cos, sin, asin, sqrt
import numpy as np
from km_matcher import KMMatcher
import scipy.stats as stats


from pandas.tseries.offsets import BDay




# Implementated based on the PyTorch
# import the resource allocation function
# replace it with your algorithm when applying LyDROO in other problems
from ResourceAllocation_3_queue import Algo1_NUM

import math

# five months
# three VM type: small, medium and large
# 0.5-0.6, 0.6-0.7, 0.7-0.8, 0.8-0.9, 0.9-1.0
# 1.5-1.6, 1.6-1.7,1.7-1.8,1.8-1.9,1.9-2.0

# small: 1.2 -2.5
# medi: 4 -7
# large: 14 -16



deal_para=0.5
N =10                     # 用户数量
M=4                  # VM types
D=4                 # datacenter number
J=1               # 总共的服务器节点
utili = np.zeros((D, M))

# 每个区域虚拟机数量的最大值
Num_max = 8
p_min_s = np.zeros((D, M))  # 算力区域池代理估计的最小值
p_max_s = np.zeros((D, M))  # 算力区域池代理估计的最大值
u_bid_s = np.zeros((D, M))  # 用户估计的最大值
u_best_s = np.zeros((D, M))  # 用户估计的最大值
p_Off_s = np.zeros((D, M))
p_best_s = np.zeros((D, M))  # 用户估计的最大值
p_final = np.zeros((D, M))  # 最终的交易值
pro_tru = np.zeros((D, J, M))  # 提供商的真实出价
user_tru = np.zeros((N))  # 用户的真实出价

User_request = pd.read_csv('submit.csv')
VM_price = pd.read_csv('VM_price.csv')
# 虚拟机类型的位置,删减相同元素，按照原始顺序进行排序
P_stan = np.array(list(VM_price.loc[:, 'price']))
P_stan = P_stan.reshape(D, M)
sigma = 0.2


def haversine(lon1, lat1, lon2, lat2): # 经度1，纬度1，经度2，纬度2 （十进制度数）
    """
    Calculate the great circle distance between two points
    on the earth (specified in decimal degrees)
    """
    # 将十进制度数转化为弧度
    lon1, lat1, lon2, lat2 = map(radians, [lon1, lat1, lon2, lat2])

    # haversine公式
    dlon = lon2 - lon1
    dlat = lat2 - lat1
    a = sin(dlat/2)**2 + cos(lat1) * cos(lat2) * sin(dlon/2)**2
    c = 2 * asin(sqrt(a))
    r = 6371 # 地球平均半径，单位为公里
    return c * r

def VM_read():
    # 8个特征，4个VM类型，5个区域
    VM_Pra=[[1,0,0,0,0,2,4,1.5],
            [1, 0, 0, 0, 0, 8, 16, 3],
            [1, 0, 0, 0, 0, 16, 32, 6],
            [1, 0, 0, 0, 0, 32, 64, 12],
            [0,1,0,0,0,2,4,1.5],
            [0, 1, 0, 0, 0, 8, 16, 3],
            [0, 1, 0, 0, 0, 16, 32, 6],
            [0, 1, 0, 0, 0, 32, 64, 12],
            [0,0,1,0,0,2,4,1.5],
            [0, 0, 1, 0, 0, 8, 16, 3],
            [0, 0, 1, 0, 0, 16, 32, 6],
            [0, 0, 1, 0, 0, 32, 64, 12],
            [0,0,0,1,0,2,4,1.5],
            [0, 0, 0, 0, 1, 8, 16, 3],
            [0, 0, 0, 0, 1, 16, 32, 6],
            [0, 0, 0, 0, 1, 32, 64, 12],
            ]
    d0 = np.vstack((VM_Pra[0:4], VM_Pra[0:4], VM_Pra[0:4], VM_Pra[0:4], VM_Pra[0:4], VM_Pra[0:4], VM_Pra[0:4],
                    VM_Pra[0:4], VM_Pra[0:4], VM_Pra[0:4]))
    d1 = np.vstack((VM_Pra[4:8], VM_Pra[4:8], VM_Pra[4:8], VM_Pra[4:8], VM_Pra[4:8], VM_Pra[4:8], VM_Pra[4:8],
                    VM_Pra[4:8], VM_Pra[4:8], VM_Pra[4:8]))
    d2 = np.vstack((VM_Pra[8:12], VM_Pra[8:12], VM_Pra[8:12], VM_Pra[8:12], VM_Pra[8:12], VM_Pra[8:12], VM_Pra[8:12],
                    VM_Pra[8:12], VM_Pra[8:12], VM_Pra[8:12]))
    d3 = np.vstack((VM_Pra[12:16], VM_Pra[12:16], VM_Pra[12:16], VM_Pra[12:16], VM_Pra[12:16], VM_Pra[12:16],
                    VM_Pra[12:16], VM_Pra[12:16], VM_Pra[12:16], VM_Pra[12:16]))
    VM_input=np.vstack((d0,d1,d2,d3))

    return VM_input

def work_read(i_idx,N):
    work_now=[]
    work = list(User_request.loc[:, 'bytes'])
    # print((i_idx+1)*N-i_idx*N)
    for i in range(i_idx*N,(i_idx+1)*N):
        work_now.append(work[i])
    work_now=np.array(work_now)
    return work_now

def dis_read():
    dis = pd.read_csv('distance.csv')
    distance=np.array(dis)
    return distance

def Unchange_price():
    VM_price = pd.read_csv('VM_price.csv')
# 虚拟机类型的位置,删减相同元素，按照原始顺序进行排序
    P_stan = np.array(list(VM_price.loc[:, 'price']))
    P_stan = P_stan.reshape(D, M)
    return P_stan




# def distance(i_idx,N,d):
#     dist=[]
#     # 用户的位置
#     U_lat = list(User_request.loc[:, 'latitude'])
#     U_lon = list(User_request.loc[:, 'longitude'])
#     # 虚拟机类型的位置,删减相同元素，按照原始顺序进行排序
#     R_lat_old = list(VM_price.loc[:, 'latitude'])
#     R_lon_old = list(VM_price.loc[:, 'longitude'])
#     R_lat=list(set(R_lat_old))
#     R_lat.sort(key = R_lat_old.index)
#     R_lon=list(set(R_lon_old ))
#     R_lon.sort(key = R_lon_old.index)
#     # 从原始的数据表中按顺序取值
#     for i in range(i_idx*N,(i_idx+1)*N):
#         dist.append(haversine(U_lon[i],U_lat[i],R_lon[d],R_lat[d]))
#
#     dist=np.array(dist)
#     return dist

def market_price(i_idx,D,M,SLA_time,VMnum_total,Num_pro):

    for d in range(D):
        for k in range(M):
            utili[d,k]=VMnum_total[d,k]/(Num_max*Num_pro)
            # print('utili[d,k]=',utili[d,k])
            # 随机用户对于虚拟机类型k的最低估价
            # p_min_s[d,k]=P_stan[d,k]-0.1*abs(sin(i_idx+1))
            # p_max_s[d,k]=P_stan[d,k]+0.1*abs(sin(i_idx+1))
            p_min_s[d, k] = P_stan[d, k] - 1
            p_max_s[d, k] = P_stan[d, k] + 2
            # print( p_min_s[d,k],p_max_s[d,k])
            # 用户的估计和最佳出价
            # u_bid_s[d,k],u_best_s[d,k]=user_evaluate(p_min_s[d,k],p_max_s[d,k])
            u_bid_s[d, k], u_best_s[d, k] = user_evaluate(p_min_s[d, k], p_max_s[d, k], SLA_time)
            # 提供商的估计和最佳要价
            p_Off_s[d,k],p_best_s[d,k]=provider_evaluate(p_min_s[d,k],p_max_s[d,k],utili[d,k])
            # 最终的市场指导成交价
            p_final[d,k]=deal_para*u_best_s[d,k]+(1-deal_para)*p_best_s[d,k]

            # # 生成实际的价格
            # pro_tru[d, :, k] = (
            #     stats.truncnorm((p_min_s[d, k] * 2 - P_stan[d,k]) / sigma, (p_max_s[d, k] * 2 - P_stan[d,k]) / sigma, loc= P_stan[d,k],
            #                     scale=sigma)).rvs(J)  # 虚拟机的价格估计


    return p_final,u_bid_s,p_max_s,p_min_s,p_Off_s

def user_evaluate(a,b,SLA_time):
    # f_v=k_v+(1-k_v)*(1/d)**delta  # between 0-1
    # user_eval=r_min+f_v*(r_max-r_min)
    # user_eval=random.uniform(a,b)
    k_v = 0.2
    f_v = k_v + (1 - k_v) * (SLA_time) ** 0.01  # between 0-1
    user_eval = a + (b - a) * f_v
    user_best=1/(deal_para+1)*user_eval+deal_para/(2*deal_para+1)*a+deal_para**2/((2*deal_para+1)*(deal_para+1))*b
    return user_eval,user_best

def provider_evaluate(a,b,utili):
    provider_eval=a+(b-a)*utili
    provider_best=1/(deal_para+1)*provider_eval+deal_para/(2*deal_para+1)*b+deal_para**2/((2*deal_para+1)*(deal_para+1))*a
    return provider_eval,provider_best

# def auction_match(u_bid_actual,freque,VMM_num):
#     # 计算权重
#     weights=np.zeros((D,N,J,M))
#     freque_max=np.ones((D,J,M))
#     cost_factor=0.4
#     for i in range(D):
#         for k in range(M):
#             for j in range(J):
#                 for n in range(N):
#                     weights[i,n,j,k]=u_bid_actual[i,n,j,k]-cost_factor*(freque_max[i,j,k]-VMM_num[i,n,k]*freque[i,n,k])**3
#         matcher = KMMatcher(weights[i,:,:,k])
#         print('weights[i,:,:,k]',np.shape(weights[i,:,:,k]))
#         best = matcher.solve(verbose=True)
#     # print( x, y,weights,sum)
#     return  best




    # print('p_Offp_Off',p_Off)
    #
    # for t in range(T):
    #     # 满足second最高价，third最低价
    #     if p_final[0,t,2]<max(p_final[0,t,:]):
    #         a_0=np.argmax(p_final[0,t,:])
    #         aa_0=p_final[0,t,2]
    #         p_final[0,t,2]=p_final[0,t,a_0]
    #         p_final[0,t,a_0]=aa_0
    #     if p_final[0,t,3]>min(p_final[0,t,:]):
    #         b_0=np.argmin(p_final[0,t,:])
    #         bb_0=p_final[0,t,3]
    #         p_final[0,t,3]=p_final[0,t,b_0]
    #         p_final[0,t,b_0]=bb_0
        ## 第二个类型
        # if p_final[1,t,2]<max(p_final[1,t,:]):
        #     a_1=np.argmax(p_final[1,t,:])
        #     aa_1=p_final[1,t,2]
        #     p_final[1,t,2]=p_final[1,t,a_1]
        #     p_final[1,t,a_1]=aa_1
        # if p_final[1,t,3]>min(p_final[1,t,:]):
        #     b_1=np.argmin(p_final[1,t,:])
        #     bb_1=p_final[1,t,3]
        #     p_final[1,t,3]=p_final[1,t,b_1]
        #     p_final[1,t,b_1]=bb_1
        # ## 第三个类型
        # if p_final[2,t,2]<max(p_final[2,t,:]):
        #     a_2=np.argmax(p_final[2,t,:])
        #     aa_2=p_final[2,t,2]
        #     p_final[2,t,2]=p_final[2,t,a_2]
        #     p_final[2,t,a_2]=aa_2
        # if p_final[2,t,3]>min(p_final[2,t,:]):
        #     b_2=np.argmin(p_final[2,t,:])
        #     bb_2=p_final[2,t,3]
        #     p_final[2,t,3]=p_final[2,t,b_2]
        #     p_final[2,t,b_2]=bb_2

    # print('p_finalp_final(d)',p_final[0 ,1:5,:])
    # df = pd.DataFrame({'Math Admin Date':pd.date_range(start=dt.datetime(2021,8,1), end = dt.datetime(2021,12,31))})
    # df.to_csv('p_final_s_1.csv')
    # # pp=p_final[1]
    # datas = pd.DataFrame(p_final[0])
    # datas.to_csv('p_final_s_2.csv')
    # df1 = pd.read_csv("p_final_s_1.csv")
    # df2 = pd.read_csv("p_final_s_2.csv")
    # df = pd.merge(df1, df2)
    # df.drop_duplicates()  #数据去重
    # df.to_csv('p_final_s.csv',encoding = 'utf-8')


