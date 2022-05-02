% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% INTNUM 0
	addi	r12, r0, 0
	sw	t0(r0), r12
	
	% INTNUM 2
	addi	r12, r0, 2
	sw	t1(r0), r12
	muli	r12, r12, 0
	lw	r11, t0(r0)
	muli	r10, r11, 4
	add	r12, r12, r10
	lw	r11, t1(r0)
	muli	r10, r11, 4
	add	r12, r12, r10
	
	% INTNUM 10
	addi	r10, r0, 10
	sw	t2(r0), r10
	
	% ASSIGN arr(r12) = t2
	lw	r10, t2(r0)
	sw	arr(r12), r10
	
	% INTNUM 1
	addi	r10, r0, 1
	sw	t3(r0), r10
	
	% INTNUM 1
	addi	r10, r0, 1
	sw	t4(r0), r10
	muli	r10, r10, 0
	lw	r12, t3(r0)
	muli	r11, r12, 4
	add	r10, r10, r11
	lw	r12, t4(r0)
	muli	r11, r12, 4
	add	r10, r10, r11
	
	% INTNUM 10
	addi	r11, r0, 10
	sw	t5(r0), r11
	
	% ASSIGN arr(r10) = t5
	lw	r11, t5(r0)
	sw	arr(r10), r11
	
	% INTNUM 0
	addi	r11, r0, 0
	sw	t6(r0), r11
	
	% INTNUM 1
	addi	r11, r0, 1
	sw	t7(r0), r11
	
	% INTNUM 2
	addi	r11, r0, 2
	sw	t8(r0), r11
	
	% MULT
	lw	r10, t7(r0)
	lw	r12, t8(r0)
	mul	r11, r10, r12
	sw	t9(r0), r11
	
	% INTNUM 1
	addi	r11, r0, 1
	sw	t10(r0), r11
	
	% PLUS
	lw	r10, t9(r0)
	lw	r12, t10(r0)
	add	r11, r10, r12
	sw	t11(r0), r11
	
	% INTNUM 1
	addi	r11, r0, 1
	sw	t12(r0), r11
	
	% INTNUM 1
	addi	r11, r0, 1
	sw	t13(r0), r11
	
	% DIV
	lw	r10, t12(r0)
	lw	r12, t13(r0)
	div	r11, r10, r12
	sw	t14(r0), r11
	
	% SUB
	lw	r10, t11(r0)
	lw	r12, t14(r0)
	sub	r11, r10, r12
	sw	t15(r0), r11
	muli	r11, r11, 0
	lw	r10, t6(r0)
	muli	r12, r10, 4
	add	r11, r11, r12
	lw	r10, t15(r0)
	muli	r12, r10, 4
	add	r11, r11, r12
	
	% INTNUM 1
	addi	r12, r0, 1
	sw	t16(r0), r12
	
	% INTNUM 1
	addi	r12, r0, 1
	sw	t17(r0), r12
	muli	r12, r12, 0
	lw	r10, t16(r0)
	muli	r9, r10, 4
	add	r12, r12, r9
	lw	r10, t17(r0)
	muli	r9, r10, 4
	add	r12, r12, r9
	
	% MULT
	lw	r10, arr(r11)
	lw	r8, arr(r12)
	mul	r9, r10, r8
	sw	t18(r0), r9
	
	% ASSIGN x = t18
	lw	r9, t18(r0)
	sw	x(r0), r9
	
	% WRITE(x)
	lw	r9, x(r0)
	sw	-8(r14), r9	% intstr arg1
	addi	r9, r0, wbuf
	sw	-12(r14), r9	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	hlt

% Data:
arr	res	8		% Space for variable arr
x	res	4		% Space for variable x
t0	res	4		% Space for variable t0
t1	res	4		% Space for variable t1
t2	res	4		% Space for variable t2
t3	res	4		% Space for variable t3
t4	res	4		% Space for variable t4
t5	res	4		% Space for variable t5
t6	res	4		% Space for variable t6
t7	res	4		% Space for variable t7
t8	res	4		% Space for variable t8
t9	res	4		% Space for variable t9
t10	res	4		% Space for variable t10
t11	res	4		% Space for variable t11
t12	res	4		% Space for variable t12
t13	res	4		% Space for variable t13
t14	res	4		% Space for variable t14
t15	res	4		% Space for variable t15
t16	res	4		% Space for variable t16
t17	res	4		% Space for variable t17
t18	res	4		% Space for variable t18
wbuf	res	32		% Buffer for printing

