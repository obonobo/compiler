% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% INTNUM 1
	addi	r12, r0, 1
	sw	t0(r0), r12
	
	% INTNUM 5
	addi	r12, r0, 5
	sw	t1(r0), r12
	
	% PLUS
	lw	r11, t0(r0)
	lw	r10, t1(r0)
	add	r12, r11, r10
	sw	t2(r0), r12
	
	% INTNUM 3
	addi	r12, r0, 3
	sw	t3(r0), r12
	
	% SUB
	lw	r11, t2(r0)
	lw	r10, t3(r0)
	sub	r12, r11, r10
	sw	t4(r0), r12
	
	% WRITE(t4)
	lw	r12, t4(r0)
	sw	-8(r14), r12	% intstr arg1
	addi	r12, r0, buf
	sw	-12(r14), r12	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	hlt

% Data:
x	res	4		% Space for variable x
y	res	4		% Space for variable y
t0	res	4		% Space for variable t0
t1	res	4		% Space for variable t1
t2	res	4		% Space for variable t2
t3	res	4		% Space for variable t3
t4	res	4		% Space for variable t4
buf	res	32		% Buffer for printing

