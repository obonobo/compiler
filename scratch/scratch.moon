% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% INTNUM 1
	addi	r12, r0, 1
	sw	t0(r0), r12
	muli	r12, r12, 0
	lw	r11, t0(r0)
	muli	r10, r11, 4
	add	r12, r12, r10
	
	% ASSIGN x = arr(r12)
	lw	r10, arr(r12)
	sw	x(r0), r10
	
	% WRITE(x)
	lw	r10, x(r0)
	sw	-8(r14), r10	% intstr arg1
	addi	r10, r0, wbuf
	sw	-12(r14), r10	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	hlt

% Data:
arr	res	8		% Space for variable arr
x	res	4		% Space for variable x
t0	res	4		% Space for variable t0
wbuf	res	32		% Buffer for printing

