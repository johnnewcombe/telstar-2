1Y�!�* . "�!Y"���UNLIST Version 1.1

 :] � �]��UNLIST modifies a BBC BASIC program so that it can be
RUN but cannot be LISTed or CONVERTed to ASCII.  This
modification cannot be reversed!
Its options are as follows:

UNLIST               Prints this description.
UNLIST filename      Modifies filename.BBC so that it
                     cannot be LISTed.  The original
                     file is renamed filename.BAK.

Note that the same restrictions apply to UNLIST as to
RENUMBER; that is, the program must not use calculated
line numbers in GOTO, GOSUB or RESTORE statements. ��!\ $ ���!~ ��e !x ����x�(S�G�x�x!���xw#���H >��� �����́*�s#r#"��[���R�8���Insufficient memory ǯ*�w#w2h 2| ��0�x�(%�͞�xO�xG͗y͞x͞��G�x͞�կ͞>�͞>�͞ͽͭe !{ ���Wl !\  ��e !x ��u !{ ���^e !~ ��u !x ���^�BBCBAK$$$�~ W����_z����W�!Y^#V#z�(��B� ��  ��>�m<���
Close error ��e>�m��������End of file - abort ͭ��e>�m����Disk full �\ >�m�6  ���<���File not found ��>�m>�m�6 �<���Directory full �\ >\ >>��m�����O� ������F�����\ ! �N ��	~˹�q ������� �N�!��	w���˹�q�����~��>͞������~#������_>�m��