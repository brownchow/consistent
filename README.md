# Go语言实现一致性哈希算法



## 解决的问题：

https://www.bilibili.com/video/BV1Hs411j73w

一致性哈希算法主要是用来解决分布式缓存中进行节点扩容（负载均衡场景），增加新的物理节点时，避免移动所有的数据，没有一致性哈希算法之前，最简单的哈希算法是：**将key进行Hash，得到一个数值，然后对物理节点的数量取模，得到的数就可以确定数据存放在第几个物理节点。**

问题是：新增物理节点时，需要对所有的key进行Hash运算，再对新的节点总数进行取模，来确定新的存放位置。这样，**大量**原来的 key 位置就会变化。有没有办法在新增节点后，**减少 key 的变动**，就是一致性哈希算法需要解决的问题

目标就是避免所有数据移动，移动尽量少的数据



## 一些问题

1、什么是哈希：

任意长度的输入，通过哈希算法，变成固定长度的输出，那个输出值，可以叫哈希值。其实哈希就是一种压缩映射

2、哈希的特点：

只要输入相同，输出就相同

3、哈希表是什么：

不准确，但是往往表达了哈希的整体意图

4、哈希碰撞：

`x != y   f(x) = f(y) `   山大王晓云， 谷歌 SHA1碰撞？

5、哈希的用途：

Redis，散列算法，分布式数据库，分布式事务，理论上只要设计到分布式，哈希是逃不掉的

6、常用的哈希函数：

直接寻址

数字分析

平方取中

折中

随机 twitter  idworker

余数

推论：其实你也经常会写哈希算法。ID 生成，文件不变（不一定使用MD5），密码加密，验证码生成

7、为什么分布式必须得有所谓的哈希算法：

涉及到一个问题：哈希一致性算法是怎么来的，为什么需要哈希一致性算法。其实也是雪崩的一个原因，且常见的，非常重要的一个原因，70%的原因

8、什么是缓存雪崩？

**尽可能让某个key（哈希值）保存在原有的服务器上**，进而使访问不要出错



## TODO:

### 1、实现 FNV 哈希算法

FNV 是个基本的，简单的，理论上背过可以可以手码的

### 2、实现PAXOS算法

### 3、查阅 Nginx中一致性哈希算法实现

https://www.youtube.com/watch?v=zvmFKwtGaYs

### 4、优化

https://github.com/kkdai/consistent/issues/2