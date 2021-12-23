# import random
# def partition(A, start, end):
#     pivot = A[end]
#     i = start - 1
#     for j in range(start, end):
#         if A[j] <= pivot:
#             i += 1
#             A[i], A[j] = A[j], A[i] 
#     A[i+1], A[end] = A[end], A[i+1]
#     print(A)
#     return i+1
        
# def quicksort(A, start, end):
#     if start < end:
#         pivot_position = partition(A, start, end)
#         quicksort(A, start, pivot_position -1)
#         quicksort(A,pivot_position + 1, end)
    
        
# A = list(range(1,11))
# random.shuffle(A)
# print(A)
# quicksort(A, 0, len(A)-1)
# print(A)