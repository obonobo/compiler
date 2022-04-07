% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% INTNUM 0
	addi	r12, r0, 0
	sw	t0(r0), r12
	
	% ASSIGN i = t0
	lw	r12, t0(r0)
	sw	i(r0), r12
	
	% While statement start
dowhile1	nop
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t2(r0), r12
	
	% Less than <
	lw	r11, i(r0)
	lw	r10, t2(r0)
	clt	r12, r11, r10
	sw	t3(r0), r12
	lw	r12, t3(r0)
	bz	r12, endwhile1
	
	% WRITE(i)
	lw	r11, i(r0)
	sw	-8(r14), r11	% intstr arg1
	addi	r11, r0, wbuf
	sw	-12(r14), r11	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	
	% INTNUM 1
	addi	r11, r0, 1
	sw	t4(r0), r11
	
	% PLUS
	lw	r10, i(r0)
	lw	r9, t4(r0)
	add	r11, r10, r9
	sw	t5(r0), r11
	
	% ASSIGN i = t5
	lw	r11, t5(r0)
	sw	i(r0), r11
	j	dowhile1
endwhile1	nop
	
	% While statement end
	hlt

% Data:
i	res	4		% Space for variable i
t0	res	4		% Space for variable t0
t2	res	4		% Space for variable t2
t3	res	4		% Space for variable t3
wbuf	res	32		% Buffer for printing
t4	res	4		% Space for variable t4
t5	res	4		% Space for variable t5

