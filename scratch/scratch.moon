% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t0(r0), r12
	
	% ASSIGN x = t0
	lw	r12, t0(r0)
	sw	x(r0), r12
	
	% INTNUM 30
	addi	r12, r0, 30
	sw	t1(r0), r12
	
	% ASSIGN y = t1
	lw	r12, t1(r0)
	sw	y(r0), r12
	
	% INTNUM 5
	addi	r12, r0, 5
	sw	t2(r0), r12
	
	% ASSIGN z = t2
	lw	r12, t2(r0)
	sw	z(r0), r12
	
	% INTNUM 4
	addi	r12, r0, 4
	sw	t3(r0), r12
	
	% ASSIGN p = t3
	lw	r12, t3(r0)
	sw	p(r0), r12
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t4(r0), r12
	
	% INTNUM 30
	addi	r12, r0, 30
	sw	t5(r0), r12
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t6(r0), r12
	
	% DIV
	lw	r11, t5(r0)
	lw	r10, t6(r0)
	div	r12, r11, r10
	sw	t7(r0), r12
	
	% PLUS
	lw	r11, t4(r0)
	lw	r10, t7(r0)
	add	r12, r11, r10
	sw	t8(r0), r12
	
	% ASSIGN q = t8
	lw	r12, t8(r0)
	sw	q(r0), r12
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t9(r0), r12
	
	% INTNUM 5
	addi	r12, r0, 5
	sw	t10(r0), r12
	
	% INTNUM 30
	addi	r12, r0, 30
	sw	t11(r0), r12
	
	% MULT
	lw	r11, t10(r0)
	lw	r10, t11(r0)
	mul	r12, r11, r10
	sw	t12(r0), r12
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t13(r0), r12
	
	% DIV
	lw	r11, t12(r0)
	lw	r10, t13(r0)
	div	r12, r11, r10
	sw	t14(r0), r12
	
	% PLUS
	lw	r11, t9(r0)
	lw	r10, t14(r0)
	add	r12, r11, r10
	sw	t15(r0), r12
	
	% ASSIGN rr = t15
	lw	r12, t15(r0)
	sw	rr(r0), r12
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t16(r0), r12
	
	% INTNUM 5
	addi	r12, r0, 5
	sw	t17(r0), r12
	
	% INTNUM 30
	addi	r12, r0, 30
	sw	t18(r0), r12
	
	% MULT
	lw	r11, t17(r0)
	lw	r10, t18(r0)
	mul	r12, r11, r10
	sw	t19(r0), r12
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t20(r0), r12
	
	% DIV
	lw	r11, t19(r0)
	lw	r10, t20(r0)
	div	r12, r11, r10
	sw	t21(r0), r12
	
	% PLUS
	lw	r11, t16(r0)
	lw	r10, t21(r0)
	add	r12, r11, r10
	sw	t22(r0), r12
	
	% INTNUM 4
	addi	r12, r0, 4
	sw	t23(r0), r12
	
	% SUB
	lw	r11, t22(r0)
	lw	r10, t23(r0)
	sub	r12, r11, r10
	sw	t24(r0), r12
	
	% ASSIGN s = t24
	lw	r12, t24(r0)
	sw	s(r0), r12
	
	% WRITE(q)
	lw	r12, q(r0)
	sw	-8(r14), r12	% intstr arg1
	addi	r12, r0, wbuf
	sw	-12(r14), r12	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	
	% WRITE(rr)
	lw	r12, rr(r0)
	sw	-8(r14), r12	% intstr arg1
	addi	r12, r0, wbuf
	sw	-12(r14), r12	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	
	% WRITE(s)
	lw	r12, s(r0)
	sw	-8(r14), r12	% intstr arg1
	addi	r12, r0, wbuf
	sw	-12(r14), r12	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	
	% PLUS
	lw	r11, q(r0)
	lw	r10, rr(r0)
	add	r12, r11, r10
	sw	t25(r0), r12
	
	% PLUS
	lw	r11, t25(r0)
	lw	r10, s(r0)
	add	r12, r11, r10
	sw	t26(r0), r12
	
	% WRITE(t26)
	lw	r12, t26(r0)
	sw	-8(r14), r12	% intstr arg1
	addi	r12, r0, wbuf
	sw	-12(r14), r12	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	hlt

% Data:
x	res	4		% Space for variable x
y	res	4		% Space for variable y
z	res	4		% Space for variable z
p	res	4		% Space for variable p
q	res	4		% Space for variable q
rr	res	4		% Space for variable rr
s	res	4		% Space for variable s
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
t19	res	4		% Space for variable t19
t20	res	4		% Space for variable t20
t21	res	4		% Space for variable t21
t22	res	4		% Space for variable t22
t23	res	4		% Space for variable t23
t24	res	4		% Space for variable t24
wbuf	res	32		% Buffer for printing
t25	res	4		% Space for variable t25
t26	res	4		% Space for variable t26

