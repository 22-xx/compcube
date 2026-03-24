#  #################################################################
#  This file contains the main LyDROO operations, including building convolutional DNN, 
#  Storing data sample, Training DNN, and generating quantized binary offloading decisions.

#  version 1.0 -- January 2021. Written based on Tensorflow 2 
#  Liang Huang (lianghuang AT zjut.edu.cn)
#  #################################################################

from __future__ import print_function
import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers
import numpy as np
import random
import math
import matplotlib.pylab as pylab
import matplotlib.pyplot as plt


print(tf.__version__)
print(tf.keras.__version__)


# DNN network for memory
class MemoryDNN:
    def __init__(
        self,
        net,
        learning_rate = 0.01,
        training_interval=1,
        batch_size=1000,
        memory_size=10,

        output_graph=False
    ):
        self.VM_num=8
        self.net = net  # the size of the DNN
        self.training_interval = training_interval      # learn every #training_interval
        self.lr = learning_rate
        self.batch_size = batch_size
        self.memory_size = memory_size
        
        # store all binary actions
        self.enumerate_actions = []

        # stored # memory entry
        self.memory_counter = 1

        # store training cost
        self.cost_his = []

        # initialize zero memory [h, m]
        self.memory = np.zeros((self.memory_size, 40*4,11+self.VM_num))

        # construct memory network
        self._build_net()

    def _build_net(self):
        self.model = tf.keras.Sequential([
            tf.keras.layers.Conv1D(64, 3, activation='relu',input_shape=[40*4,11]),
            tf.keras.layers.Conv1D(64, 3, activation='relu'),  # second Conv1D with 32 channels and kearnal size 3
            tf.keras.layers.Conv1D(64, 3, activation='relu'),  # second Conv1D with 32 channels and kearnal size 3 
            tf.keras.layers.LSTM(128,return_sequences=True),
            tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(128)),
            tf.keras.layers.Flatten(),
            tf.keras.layers.Dense(64, activation='relu'),
            tf.keras.layers.Dense(64, activation='relu'),
            tf.keras.layers.Dense(40*4*self.VM_num, activation='sigmoid')
        ])
# 
        self.model.compile(optimizer=keras.optimizers.Adam(lr=self.lr), loss=tf.losses.binary_crossentropy, metrics=['accuracy'])

    def remember(self, h, m):
        # replace the old memory with new memory
        # 虚拟机类型
        m_new=np.zeros((40*4,self.VM_num))
        for i in range(40*4):
            m_new[i,int(m[i]-1)]=1
        idx = self.memory_counter % self.memory_size
        # 多维张量  (60,9)
        self.memory[idx, :,:] = np.hstack((h, m_new))

        self.memory_counter += 1

    def encode(self, h, m):
        # encoding the entry
        self.remember(h, m)
        # train the DNN every 10 step
#        if self.memory_counter> self.memory_size / 2 and self.memory_counter % self.training_interval == 0:
        if self.memory_counter % self.training_interval == 0:
            self.learn()

    def learn(self):
        # sample batch memory from all memory
        m_train=np.zeros((128,40*4*self.VM_num))
        if self.memory_counter > self.memory_size:
            sample_index = np.random.choice(self.memory_size, size=self.batch_size)
        else:
            sample_index = np.random.choice(self.memory_counter, size=self.batch_size)
        batch_memory = self.memory[sample_index, :,:]
        h_train = batch_memory[:,:,0:11]

        # h_train = h_train.reshape(self.batch_size,2)
        for i in range(128):
            m_train[i] = batch_memory[i, :,11:].flatten()  # (1,180)

        # print(h_train)          # (128, 10)
        # print(m_train)          # (128, 10)

        # train the DNN
        hist = self.model.fit(h_train, m_train, verbose=0)
        self.cost = hist.history['loss'][0]
        # assert(self.cost > 0)
        self.cost_his.append(self.cost)

    def decode(self, h):
        # to have batch dimension when feed into tf placeholder
        h = h[np.newaxis, :]
        m_pred = self.model.predict(h)
        m_pred = m_pred.reshape(40*4,self.VM_num)
        array_m_pred = np.array(m_pred)
        result_column = np.array([])
        for i in range(0,40*4):
            b = array_m_pred[i].argsort()[-2:][::-1]
            while(1):
                n = random.randint(0, self.VM_num)
                if n not in b:
                    b = np.append(b,n)
                    result_column = np.append(result_column,b)
                    break
        # 选取5个概率大的值，一个随机的值
        result_column = result_column.reshape(40*4,3)
        return result_column





    def plot_cost(self):
        import matplotlib.pyplot as plt

        myparams = {
            'axes.labelsize': '20',
            'xtick.labelsize': '20',
            'ytick.labelsize': '20',

            'lines.linewidth': 3.5,
            'legend.fontsize': '20',
            'font.family': 'Times New Roman'
        }
        pylab.rcParams.update(myparams)
        plt.style.use('seaborn-deep')
        fig = plt.figure(figsize=(10, 8))
        plt.plot(np.arange(len(self.cost_his))*self.training_interval, self.cost_his)
        plt.ylabel('Training Loss')
        plt.xlabel('Time Frames')
        plt.savefig(r'./figure/loss.pdf', dpi=600)
        # plt.show()
