# Development notes

## gRPC debugging

When doing a:

```
$ etcdctl --endpoints=http://localhost:2379 set /kubernetes.io/namespaces/kube-system "."
$ etcdctl get /kubernetes.io/namespaces/kube-system
```

I see the following:

```
$ sudo tcpdump -i any port 2379
06:38:12.432788 IP localhost.52588 > localhost.2379: Flags [SEW], seq 1255564419, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543945 ecr 0,sackOK,eol], length 0
06:38:12.432798 IP localhost.52588 > localhost.2379: Flags [SEW], seq 1255564419, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543945 ecr 0,sackOK,eol], length 0
06:38:12.432856 IP localhost.2379 > localhost.52588: Flags [S.], seq 3631823651, ack 1255564420, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543945 ecr 2609543945,sackOK,eol], length 0
06:38:12.432859 IP localhost.2379 > localhost.52588: Flags [S.], seq 3631823651, ack 1255564420, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543945 ecr 2609543945,sackOK,eol], length 0
06:38:12.432868 IP localhost.52588 > localhost.2379: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543945 ecr 2609543945], length 0
06:38:12.432870 IP localhost.52588 > localhost.2379: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543945 ecr 2609543945], length 0
06:38:12.432878 IP localhost.2379 > localhost.52588: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543945 ecr 2609543945], length 0
06:38:12.432879 IP localhost.2379 > localhost.52588: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543945 ecr 2609543945], length 0
06:38:12.432981 IP localhost.52588 > localhost.2379: Flags [P.], seq 1:106, ack 1, win 12759, options [nop,nop,TS val 2609543945 ecr 2609543945], length 105
06:38:12.432988 IP localhost.52588 > localhost.2379: Flags [P.], seq 1:106, ack 1, win 12759, options [nop,nop,TS val 2609543945 ecr 2609543945], length 105
06:38:12.433007 IP localhost.2379 > localhost.52588: Flags [.], ack 106, win 12756, options [nop,nop,TS val 2609543945 ecr 2609543945], length 0
06:38:12.433009 IP localhost.2379 > localhost.52588: Flags [.], ack 106, win 12756, options [nop,nop,TS val 2609543945 ecr 2609543945], length 0
06:38:12.434169 IP localhost.2379 > localhost.52588: Flags [P.], seq 1:279, ack 106, win 12756, options [nop,nop,TS val 2609543946 ecr 2609543945], length 278
06:38:12.434180 IP localhost.2379 > localhost.52588: Flags [P.], seq 1:279, ack 106, win 12756, options [nop,nop,TS val 2609543946 ecr 2609543945], length 278
06:38:12.434192 IP localhost.52588 > localhost.2379: Flags [.], ack 279, win 12750, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434195 IP localhost.52588 > localhost.2379: Flags [.], ack 279, win 12750, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434549 IP localhost.52589 > localhost.2379: Flags [S], seq 3305837192, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543946 ecr 0,sackOK,eol], length 0
06:38:12.434558 IP localhost.52589 > localhost.2379: Flags [S], seq 3305837192, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543946 ecr 0,sackOK,eol], length 0
06:38:12.434604 IP localhost.2379 > localhost.52589: Flags [S.], seq 2229774752, ack 3305837193, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543946 ecr 2609543946,sackOK,eol], length 0
06:38:12.434608 IP localhost.2379 > localhost.52589: Flags [S.], seq 2229774752, ack 3305837193, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609543946 ecr 2609543946,sackOK,eol], length 0
06:38:12.434614 IP localhost.52589 > localhost.2379: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434615 IP localhost.52589 > localhost.2379: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434621 IP localhost.2379 > localhost.52589: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434622 IP localhost.2379 > localhost.52589: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434684 IP localhost.52589 > localhost.2379: Flags [P.], seq 1:180, ack 1, win 12759, options [nop,nop,TS val 2609543946 ecr 2609543946], length 179
06:38:12.434692 IP localhost.52589 > localhost.2379: Flags [P.], seq 1:180, ack 1, win 12759, options [nop,nop,TS val 2609543946 ecr 2609543946], length 179
06:38:12.434700 IP localhost.2379 > localhost.52589: Flags [.], ack 180, win 12753, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.434702 IP localhost.2379 > localhost.52589: Flags [.], ack 180, win 12753, options [nop,nop,TS val 2609543946 ecr 2609543946], length 0
06:38:12.435588 IP localhost.2379 > localhost.52589: Flags [P.], seq 1:316, ack 180, win 12753, options [nop,nop,TS val 2609543947 ecr 2609543946], length 315
06:38:12.435598 IP localhost.2379 > localhost.52589: Flags [P.], seq 1:316, ack 180, win 12753, options [nop,nop,TS val 2609543947 ecr 2609543946], length 315
06:38:12.435607 IP localhost.52589 > localhost.2379: Flags [.], ack 316, win 12749, options [nop,nop,TS val 2609543947 ecr 2609543947], length 0
06:38:12.435609 IP localhost.52589 > localhost.2379: Flags [.], ack 316, win 12749, options [nop,nop,TS val 2609543947 ecr 2609543947], length 0
06:38:12.436981 IP localhost.52589 > localhost.2379: Flags [F.], seq 180, ack 316, win 12749, options [nop,nop,TS val 2609543948 ecr 2609543947], length 0
06:38:12.436994 IP localhost.52589 > localhost.2379: Flags [F.], seq 180, ack 316, win 12749, options [nop,nop,TS val 2609543948 ecr 2609543947], length 0
06:38:12.436998 IP localhost.52588 > localhost.2379: Flags [F.], seq 106, ack 279, win 12750, options [nop,nop,TS val 2609543948 ecr 2609543946], length 0
06:38:12.437005 IP localhost.2379 > localhost.52589: Flags [.], ack 181, win 12753, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437008 IP localhost.52588 > localhost.2379: Flags [F.], seq 106, ack 279, win 12750, options [nop,nop,TS val 2609543948 ecr 2609543946], length 0
06:38:12.437008 IP localhost.2379 > localhost.52589: Flags [.], ack 181, win 12753, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437012 IP localhost.2379 > localhost.52588: Flags [.], ack 107, win 12756, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437016 IP localhost.2379 > localhost.52588: Flags [.], ack 107, win 12756, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437639 IP localhost.2379 > localhost.52588: Flags [F.], seq 279, ack 107, win 12756, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437650 IP localhost.2379 > localhost.52588: Flags [F.], seq 279, ack 107, win 12756, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437668 IP localhost.2379 > localhost.52589: Flags [F.], seq 316, ack 181, win 12753, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437675 IP localhost.52588 > localhost.2379: Flags [.], ack 280, win 12750, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437678 IP localhost.2379 > localhost.52589: Flags [F.], seq 316, ack 181, win 12753, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437680 IP localhost.52588 > localhost.2379: Flags [.], ack 280, win 12750, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437695 IP localhost.52589 > localhost.2379: Flags [.], ack 317, win 12749, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0
06:38:12.437703 IP localhost.52589 > localhost.2379: Flags [.], ack 317, win 12749, options [nop,nop,TS val 2609543948 ecr 2609543948], length 0

```

Running `test-e2e-etcd3.sh` produces:

```
06:40:54.792278 IP6 localhost.52616 > localhost.2379: Flags [SEW], seq 2089114608, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609704804 ecr 0,sackOK,eol], length 0
06:40:54.792286 IP6 localhost.52616 > localhost.2379: Flags [SEW], seq 2089114608, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609704804 ecr 0,sackOK,eol], length 0
06:40:54.792353 IP6 localhost.2379 > localhost.52616: Flags [S.E], seq 1050182938, ack 2089114609, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609704804 ecr 2609704804,sackOK,eol], length 0
06:40:54.792356 IP6 localhost.2379 > localhost.52616: Flags [S.E], seq 1050182938, ack 2089114609, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609704804 ecr 2609704804,sackOK,eol], length 0
06:40:54.792364 IP6 localhost.52616 > localhost.2379: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609704804 ecr 2609704804], length 0
06:40:54.792365 IP6 localhost.52616 > localhost.2379: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609704804 ecr 2609704804], length 0
06:40:54.792373 IP6 localhost.2379 > localhost.52616: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609704804 ecr 2609704804], length 0
06:40:54.792375 IP6 localhost.2379 > localhost.52616: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609704804 ecr 2609704804], length 0
06:40:54.792905 IP6 localhost.52616 > localhost.2379: Flags [P.], seq 1:106, ack 1, win 12743, options [nop,nop,TS val 2609704804 ecr 2609704804], length 105
06:40:54.792922 IP6 localhost.52616 > localhost.2379: Flags [P.], seq 1:106, ack 1, win 12743, options [nop,nop,TS val 2609704804 ecr 2609704804], length 105
06:40:54.792938 IP6 localhost.2379 > localhost.52616: Flags [.], ack 106, win 12740, options [nop,nop,TS val 2609704804 ecr 2609704804], length 0
06:40:54.792944 IP6 localhost.2379 > localhost.52616: Flags [.], ack 106, win 12740, options [nop,nop,TS val 2609704804 ecr 2609704804], length 0
06:40:54.793754 IP6 localhost.2379 > localhost.52616: Flags [P.], seq 1:279, ack 106, win 12740, options [nop,nop,TS val 2609704805 ecr 2609704804], length 278
06:40:54.793766 IP6 localhost.2379 > localhost.52616: Flags [P.], seq 1:279, ack 106, win 12740, options [nop,nop,TS val 2609704805 ecr 2609704804], length 278
06:40:54.793785 IP6 localhost.52616 > localhost.2379: Flags [.], ack 279, win 12735, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.793790 IP6 localhost.52616 > localhost.2379: Flags [.], ack 279, win 12735, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.794337 IP localhost.52617 > localhost.2379: Flags [SEW], seq 1486709102, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609704805 ecr 0,sackOK,eol], length 0
06:40:54.794346 IP localhost.52617 > localhost.2379: Flags [SEW], seq 1486709102, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609704805 ecr 0,sackOK,eol], length 0
06:40:54.794412 IP localhost.2379 > localhost.52617: Flags [S.], seq 2700600856, ack 1486709103, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609704805 ecr 2609704805,sackOK,eol], length 0
06:40:54.794415 IP localhost.2379 > localhost.52617: Flags [S.], seq 2700600856, ack 1486709103, win 65535, options [mss 16344,nop,wscale 5,nop,nop,TS val 2609704805 ecr 2609704805,sackOK,eol], length 0
06:40:54.794423 IP localhost.52617 > localhost.2379: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.794424 IP localhost.52617 > localhost.2379: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.794429 IP localhost.2379 > localhost.52617: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.794431 IP localhost.2379 > localhost.52617: Flags [.], ack 1, win 12759, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.794561 IP localhost.52617 > localhost.2379: Flags [P.], seq 1:213, ack 1, win 12759, options [nop,nop,TS val 2609704805 ecr 2609704805], length 212
06:40:54.794570 IP localhost.52617 > localhost.2379: Flags [P.], seq 1:213, ack 1, win 12759, options [nop,nop,TS val 2609704805 ecr 2609704805], length 212
06:40:54.794586 IP localhost.2379 > localhost.52617: Flags [.], ack 213, win 12752, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.794591 IP localhost.2379 > localhost.52617: Flags [.], ack 213, win 12752, options [nop,nop,TS val 2609704805 ecr 2609704805], length 0
06:40:54.796505 IP localhost.2379 > localhost.52617: Flags [P.], seq 1:321, ack 213, win 12752, options [nop,nop,TS val 2609704807 ecr 2609704805], length 320
06:40:54.796516 IP localhost.2379 > localhost.52617: Flags [P.], seq 1:321, ack 213, win 12752, options [nop,nop,TS val 2609704807 ecr 2609704805], length 320
06:40:54.796528 IP localhost.52617 > localhost.2379: Flags [.], ack 321, win 12749, options [nop,nop,TS val 2609704807 ecr 2609704807], length 0
06:40:54.796531 IP localhost.52617 > localhost.2379: Flags [.], ack 321, win 12749, options [nop,nop,TS val 2609704807 ecr 2609704807], length 0
06:40:54.798935 IP localhost.52617 > localhost.2379: Flags [F.], seq 213, ack 321, win 12749, options [nop,nop,TS val 2609704809 ecr 2609704807], length 0
06:40:54.798943 IP localhost.52617 > localhost.2379: Flags [F.], seq 213, ack 321, win 12749, options [nop,nop,TS val 2609704809 ecr 2609704807], length 0
06:40:54.798949 IP6 localhost.52616 > localhost.2379: Flags [F.], seq 106, ack 279, win 12735, options [nop,nop,TS val 2609704809 ecr 2609704805], length 0
06:40:54.798960 IP localhost.2379 > localhost.52617: Flags [.], ack 214, win 12752, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.798964 IP6 localhost.52616 > localhost.2379: Flags [F.], seq 106, ack 279, win 12735, options [nop,nop,TS val 2609704809 ecr 2609704805], length 0
06:40:54.798966 IP localhost.2379 > localhost.52617: Flags [.], ack 214, win 12752, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.798976 IP6 localhost.2379 > localhost.52616: Flags [.], ack 107, win 12740, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.798984 IP6 localhost.2379 > localhost.52616: Flags [.], ack 107, win 12740, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799557 IP localhost.2379 > localhost.52617: Flags [F.], seq 321, ack 214, win 12752, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799565 IP localhost.2379 > localhost.52617: Flags [F.], seq 321, ack 214, win 12752, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799582 IP6 localhost.2379 > localhost.52616: Flags [F.], seq 279, ack 107, win 12740, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799590 IP localhost.52617 > localhost.2379: Flags [.], ack 322, win 12749, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799593 IP6 localhost.2379 > localhost.52616: Flags [F.], seq 279, ack 107, win 12740, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799595 IP localhost.52617 > localhost.2379: Flags [.], ack 322, win 12749, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799604 IP6 localhost.52616 > localhost.2379: Flags [.], ack 280, win 12735, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:40:54.799613 IP6 localhost.52616 > localhost.2379: Flags [.], ack 280, win 12735, options [nop,nop,TS val 2609704809 ecr 2609704809], length 0
06:41:00.601224 IP6 localhost.52626 > localhost.2379: Flags [S], seq 941435592, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710535 ecr 0,sackOK,eol], length 0
06:41:00.601234 IP6 localhost.52626 > localhost.2379: Flags [S], seq 941435592, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710535 ecr 0,sackOK,eol], length 0
06:41:00.601280 IP6 localhost.2379 > localhost.52626: Flags [S.], seq 2859525419, ack 941435593, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710535 ecr 2609710535,sackOK,eol], length 0
06:41:00.601285 IP6 localhost.2379 > localhost.52626: Flags [S.], seq 2859525419, ack 941435593, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710535 ecr 2609710535,sackOK,eol], length 0
06:41:00.601293 IP6 localhost.52626 > localhost.2379: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710535 ecr 2609710535], length 0
06:41:00.601294 IP6 localhost.52626 > localhost.2379: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710535 ecr 2609710535], length 0
06:41:00.601301 IP6 localhost.2379 > localhost.52626: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710535 ecr 2609710535], length 0
06:41:00.601303 IP6 localhost.2379 > localhost.52626: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710535 ecr 2609710535], length 0
06:41:00.601419 IP6 localhost.52626 > localhost.2379: Flags [P.], seq 1:103, ack 1, win 12743, options [nop,nop,TS val 2609710535 ecr 2609710535], length 102
06:41:00.601429 IP6 localhost.52626 > localhost.2379: Flags [P.], seq 1:103, ack 1, win 12743, options [nop,nop,TS val 2609710535 ecr 2609710535], length 102
06:41:00.601439 IP6 localhost.2379 > localhost.52626: Flags [.], ack 103, win 12740, options [nop,nop,TS val 2609710535 ecr 2609710535], length 0
06:41:00.601442 IP6 localhost.2379 > localhost.52626: Flags [.], ack 103, win 12740, options [nop,nop,TS val 2609710535 ecr 2609710535], length 0
06:41:00.602439 IP6 localhost.2379 > localhost.52626: Flags [P.], seq 1:153, ack 103, win 12740, options [nop,nop,TS val 2609710536 ecr 2609710535], length 152
06:41:00.602449 IP6 localhost.2379 > localhost.52626: Flags [P.], seq 1:153, ack 103, win 12740, options [nop,nop,TS val 2609710536 ecr 2609710535], length 152
06:41:00.602460 IP6 localhost.52626 > localhost.2379: Flags [.], ack 153, win 12739, options [nop,nop,TS val 2609710536 ecr 2609710536], length 0
06:41:00.602463 IP6 localhost.52626 > localhost.2379: Flags [.], ack 153, win 12739, options [nop,nop,TS val 2609710536 ecr 2609710536], length 0
06:41:00.603475 IP6 localhost.52627 > localhost.2379: Flags [SEW], seq 1709633836, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710537 ecr 0,sackOK,eol], length 0
06:41:00.603485 IP6 localhost.52627 > localhost.2379: Flags [SEW], seq 1709633836, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710537 ecr 0,sackOK,eol], length 0
06:41:00.603535 IP6 localhost.2379 > localhost.52627: Flags [S.E], seq 519321316, ack 1709633837, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710537 ecr 2609710537,sackOK,eol], length 0
06:41:00.603539 IP6 localhost.2379 > localhost.52627: Flags [S.E], seq 519321316, ack 1709633837, win 65535, options [mss 16324,nop,wscale 5,nop,nop,TS val 2609710537 ecr 2609710537,sackOK,eol], length 0
06:41:00.603550 IP6 localhost.52627 > localhost.2379: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603552 IP6 localhost.52627 > localhost.2379: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603560 IP6 localhost.2379 > localhost.52627: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603563 IP6 localhost.2379 > localhost.52627: Flags [.], ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603665 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 1:25, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 24
06:41:00.603681 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 1:25, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 24
06:41:00.603684 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 25:34, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 9
06:41:00.603696 IP6 localhost.2379 > localhost.52627: Flags [.], ack 25, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603697 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 34:47, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 13
06:41:00.603701 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 25:34, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 9
06:41:00.603702 IP6 localhost.2379 > localhost.52627: Flags [.], ack 25, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603703 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 34:47, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 13
06:41:00.603709 IP6 localhost.2379 > localhost.52627: Flags [.], ack 34, win 12742, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603719 IP6 localhost.2379 > localhost.52627: Flags [.], ack 47, win 12742, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603721 IP6 localhost.2379 > localhost.52627: Flags [.], ack 34, win 12742, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.603722 IP6 localhost.2379 > localhost.52627: Flags [.], ack 47, win 12742, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.604022 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 47:177, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 130
06:41:00.604030 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 47:177, ack 1, win 12743, options [nop,nop,TS val 2609710537 ecr 2609710537], length 130
06:41:00.604044 IP6 localhost.2379 > localhost.52627: Flags [.], ack 177, win 12738, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.604050 IP6 localhost.2379 > localhost.52627: Flags [.], ack 177, win 12738, options [nop,nop,TS val 2609710537 ecr 2609710537], length 0
06:41:00.604511 IP6 localhost.2379 > localhost.52627: Flags [P.], seq 1:23, ack 177, win 12738, options [nop,nop,TS val 2609710538 ecr 2609710537], length 22
06:41:00.604523 IP6 localhost.2379 > localhost.52627: Flags [P.], seq 1:23, ack 177, win 12738, options [nop,nop,TS val 2609710538 ecr 2609710537], length 22
06:41:00.604536 IP6 localhost.52627 > localhost.2379: Flags [.], ack 23, win 12743, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.604540 IP6 localhost.52627 > localhost.2379: Flags [.], ack 23, win 12743, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.604602 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 177:186, ack 23, win 12743, options [nop,nop,TS val 2609710538 ecr 2609710538], length 9
06:41:00.604610 IP6 localhost.52627 > localhost.2379: Flags [P.], seq 177:186, ack 23, win 12743, options [nop,nop,TS val 2609710538 ecr 2609710538], length 9
06:41:00.604621 IP6 localhost.2379 > localhost.52627: Flags [.], ack 186, win 12737, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.604623 IP6 localhost.2379 > localhost.52627: Flags [.], ack 186, win 12737, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.604866 IP6 localhost.2379 > localhost.52627: Flags [P.], seq 23:32, ack 186, win 12737, options [nop,nop,TS val 2609710538 ecr 2609710538], length 9
06:41:00.604877 IP6 localhost.2379 > localhost.52627: Flags [P.], seq 23:32, ack 186, win 12737, options [nop,nop,TS val 2609710538 ecr 2609710538], length 9
06:41:00.604889 IP6 localhost.52627 > localhost.2379: Flags [.], ack 32, win 12742, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.604892 IP6 localhost.52627 > localhost.2379: Flags [.], ack 32, win 12742, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.605085 IP6 localhost.2379 > localhost.52627: Flags [P.], seq 32:128, ack 186, win 12737, options [nop,nop,TS val 2609710538 ecr 2609710538], length 96
06:41:00.605096 IP6 localhost.2379 > localhost.52627: Flags [P.], seq 32:128, ack 186, win 12737, options [nop,nop,TS val 2609710538 ecr 2609710538], length 96
06:41:00.605109 IP6 localhost.52627 > localhost.2379: Flags [.], ack 128, win 12739, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.605113 IP6 localhost.52627 > localhost.2379: Flags [.], ack 128, win 12739, options [nop,nop,TS val 2609710538 ecr 2609710538], length 0
06:41:00.608241 IP6 localhost.52627 > localhost.2379: Flags [F.], seq 186, ack 128, win 12739, options [nop,nop,TS val 2609710541 ecr 2609710538], length 0
06:41:00.608263 IP6 localhost.2379 > localhost.52627: Flags [.], ack 187, win 12737, options [nop,nop,TS val 2609710541 ecr 2609710541], length 0
06:41:00.608270 IP6 localhost.2379 > localhost.52627: Flags [.], ack 187, win 12737, options [nop,nop,TS val 2609710541 ecr 2609710541], length 0
06:41:00.608910 IP6 localhost.2379 > localhost.52627: Flags [F.], seq 128, ack 187, win 12737, options [nop,nop,TS val 2609710541 ecr 2609710541], length 0
06:41:00.608920 IP6 localhost.2379 > localhost.52627: Flags [F.], seq 128, ack 187, win 12737, options [nop,nop,TS val 2609710541 ecr 2609710541], length 0
06:41:00.608939 IP6 localhost.52627 > localhost.2379: Flags [.], ack 129, win 12739, options [nop,nop,TS val 2609710541 ecr 2609710541], length 0
06:41:00.608942 IP6 localhost.52627 > localhost.2379: Flags [.], ack 129, win 12739, options [nop,nop,TS val 2609710541 ecr 2609710541], length 0
^[[A06:41:08.818795 IP6 localhost.52626 > localhost.2379: Flags [F.], seq 103, ack 153, win 12739, options [nop,nop,TS val 2609718653 ecr 2609710536], length 0
06:41:08.818815 IP6 localhost.52626 > localhost.2379: Flags [F.], seq 103, ack 153, win 12739, options [nop,nop,TS val 2609718653 ecr 2609710536], length 0
06:41:08.818839 IP6 localhost.2379 > localhost.52626: Flags [.], ack 104, win 12740, options [nop,nop,TS val 2609718653 ecr 2609718653], length 0
06:41:08.818844 IP6 localhost.2379 > localhost.52626: Flags [.], ack 104, win 12740, options [nop,nop,TS val 2609718653 ecr 2609718653], length 0
06:41:08.819574 IP6 localhost.2379 > localhost.52626: Flags [F.], seq 153, ack 104, win 12740, options [nop,nop,TS val 2609718653 ecr 2609718653], length 0
06:41:08.819585 IP6 localhost.2379 > localhost.52626: Flags [F.], seq 153, ack 104, win 12740, options [nop,nop,TS val 2609718653 ecr 2609718653], length 0
06:41:08.819606 IP6 localhost.52626 > localhost.2379: Flags [.], ack 154, win 12739, options [nop,nop,TS val 2609718653 ecr 2609718653], length 0
06:41:08.819610 IP6 localhost.52626 > localhost.2379: Flags [.], ack 154, win 12739, options [nop,nop,TS val 2609718653 ecr 2609718653], length 0
```

## etcd key prefixes

### Kubernetes

```
/kubernetes.io/ranges
/kubernetes.io/statefulsets
/kubernetes.io/jobs
/kubernetes.io/horizontalpodautoscalers
/kubernetes.io/events
/kubernetes.io/masterleases
/kubernetes.io/minions
/kubernetes.io/persistentvolumes
/kubernetes.io/configmaps
/kubernetes.io/controllers
/kubernetes.io/deployments
/kubernetes.io/serviceaccounts
/kubernetes.io/services
/kubernetes.io/namespaces
/kubernetes.io/securitycontextconstraints
/kubernetes.io/thirdpartyresources
/kubernetes.io/persistentvolumeclaims
/kubernetes.io/pods
/kubernetes.io/replicasets
/kubernetes.io/secrets
```

### OpenShift

```
/openshift.io/authorization
/openshift.io/buildconfigs
/openshift.io/oauth
/openshift.io/registry
/openshift.io/users
/openshift.io/useridentities
/openshift.io/builds
/openshift.io/deploymentconfigs
/openshift.io/images
/openshift.io/imagestreams
/openshift.io/ranges
/openshift.io/routes
/openshift.io/templates
```
