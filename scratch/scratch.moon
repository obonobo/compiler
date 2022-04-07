% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% If statement start
	
	% INTNUM 1
	addi	r12, r0, 1
	sw	t0(r0), r12
	
	% INTNUM 0
	addi	r12, r0, 0
	sw	t1(r0), r12
	
	% EQ(==)
	lw	r11, t0(r0)
	lw	r10, t1(r0)
	ceq	r12, r11, r10
	sw	t2(r0), r12
	lw	r12, t2(r0)
	bz	r12, else3
	
	% INTNUM 1
	addi	r11, r0, 1
	sw	t4(r0), r11
	
	% WRITE(t4)
	lw	r11, t4(r0)
	sw	-8(r14), r11	% intstr arg1
	addi	r11, r0, wbuf
	sw	-12(r14), r11	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	j	endIf3
else3	nop
	
	% INTNUM 0
	addi	r11, r0, 0
	sw	t5(r0), r11
	
	% WRITE(t5)
	lw	r11, t5(r0)
	sw	-8(r14), r11	% intstr arg1
	addi	r11, r0, wbuf
	sw	-12(r14), r11	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
endIf3	nop
	
	% If statement end
	hlt

% Data:
t0	res	4		% Space for variable t0
t1	res	4		% Space for variable t1
t2	res	4		% Space for variable t2
t4	res	4		% Space for variable t4
wbuf	res	32		% Buffer for printing
t5	res	4		% Space for variable t5

