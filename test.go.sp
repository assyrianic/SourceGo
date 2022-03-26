/**
 * file generated by the GoToSourcePawn Transpiler v1.4b
 * Copyright 2020 (C) Kevin Yonan aka Nergal, Assyrianic.
 * GoToSourcePawn Project is licensed under MIT.
 * link: 'https://github.com/assyrianic/Go2SourcePawn'
 */

#include <sourcemod>
#include <sdktools>

enum struct Point {
	float x;
	float y;
}

enum struct PlayerInfo {
	float Origin[3];
	float Angle[3];
	int Weaps[3];
	Function PutInServer;
}

enum struct ClientInfo {
	int Clients[2][66];
}


char a[] = "A";
int b = MAXPLAYERS;

char c[] = a;

char d[] = "D";

char e[] = "e1";

float f = 1.00;

char MakeStrMap[] = "StringMap smap = new StringMap();";



typedef Kektus = function Handle (const float i[3], const float x[3], const char[] b, char blocks[64], int& KC);

typedef EventFunc = function Action (Event event, const char[] name, bool dontBroadcast);

typedef VecFunc = function float (const float vec[3], float& VecFunc_param1, float& VecFunc_param2);

Plugin myself = {
	name = "SrcGo Plugin",
	author = "Nergal",
	description = "Plugin made into SP from SrcGo.",
	version = "1.0a",
	url = "https://github.com/assyrianic/Go2SourcePawn"
};

char str_array[4][] = {
	"kek",
	"foo",
	"bar",
	"bazz"
};

Function ff1;

Function ff2;

Function ff3;

public float TestOrigin(float& TestOrigin_param1, float& TestOrigin_param2)
{
	PlayerInfo pi;
	float o[3];

	return pi.GetOrigin(o, TestOrigin_param1, TestOrigin_param2);
}

native int FF1();

native int FF2();

native float FF3();

native int FF4(int& FF4_param1, float& FF4_param2);

public int GG1(int& GG1_param1, float& GG1_param2)
{
	GG1_param1 = FF2();
	GG1_param2 = FF3();
	return FF1();
}

public int GG2(int& GG2_param1, float& GG2_param2)
{
	Call_StartFunction(null, ff3);
	Call_Finish(GG2_param2);
	Call_StartFunction(null, ff2);
	Call_Finish(GG2_param1);
	int fptr_temp0;

	Call_StartFunction(null, ff1);
	Call_Finish(fptr_temp0);
	return fptr_temp0;
}

public int GG3(int& GG3_param1, float& GG3_param2)
{
	return FF4(GG3_param1, GG3_param2);
}

public void OnPluginStart()
{
	int inlined_call_res1;
	int inlined_call_res2;

	int inlined_call_res;

	Handle my_timer;

	float x;
	float y;
	float z;

	ClientInfo cinfo;

	for (int main_iter0; main_iter0 < sizeof(cinfo.Clients); main_iter0++)
	{
		int p1[66];

		p1 = cinfo.Clients[main_iter0];
		for (int main_iter1; main_iter1 < sizeof(p1); main_iter1++)
		{
			int x1;

			x1 = p1[main_iter1];
			bool is_in_game;

			is_in_game = IsClientInGame(x1);
		}
	}
	PlayerInfo p;

	float origin[3];

	x = p.GetOrigin(origin, y, z);
	int k;
	int l;

	k &= ~(l);
	Function CB = IndirectMultiRet;
	p.PutInServer = SrcGoTmpFunc0;
	for (int i = 1; i <= MaxClients; i++)
	{
		bool fptr_temp4;

		bool fptr_temp3;

		bool fptr_temp2;

		bool j;
		bool k;
		bool l;

		Call_StartFunction(null, p.PutInServer);
		Call_PushCell(i);
		Call_Finish();
		Call_StartFunction(null, CB);
		Call_PushCellRef(fptr_temp3);
		Call_PushCellRef(fptr_temp4);
		Call_Finish(fptr_temp2);
		bool fptr_temp1;

		Call_StartFunction(null, CB);
		Call_PushCellRef(k);
		Call_PushCellRef(l);
		Call_Finish(fptr_temp1);
		j = fptr_temp1;
	}
	for (float f = 2.0; f < 100.0; f = Pow(f, 2.0))
	{
		PrintToServer("%0.2f", f);
	}
	my_timer = CreateTimer(0.1, SrcGoTmpFunc1, 0, 0);
	inlined_call_res = SrcGoTmpFunc2(1, 2);
	inlined_call_res1 = SrcGoTmpFunc3(1, 2, inlined_call_res2);
	Function caller = SrcGoTmpFunc4;
	
	int n;
	Call_StartFunction(null, caller);
	Call_PushCell(1); Call_PushCell(2);
	Call_Finish(n);
	KeyValues kv;

	kv = new KeyValues("kek1", "kek_key", "kek_val");
	delete kv;
	AddMultiTargetFilter("@!party", SrcGoTmpFunc5, "The D&D Quest Party", false);
	StringMap smap = new StringMap();
	char[] new_str = new char[sizeof("100")];
	int[] new_kek = new int[inlined_call_res];
}

public bool IndirectMultiRet(bool& IndirectMultiRet_param1, bool& IndirectMultiRet_param2)
{
	return MultiRetFn(IndirectMultiRet_param1, IndirectMultiRet_param2);
}

public bool MultiRetFn(bool& MultiRetFn_param1, bool& MultiRetFn_param2)
{
	MultiRetFn_param1 = false;
	MultiRetFn_param2 = true;
	return true;
}

public void OnClientPutInServer(int client)
{
}

public void GetProjPosToScreen(int client, const float vecDelta[3], float& xpos, float& ypos)
{
	float playerAngles[3];
	float vecforward[3];
	float right[3];
	float up[3];

	GetClientEyeAngles(client, playerAngles);
	up[2] = 1.0;
	GetAngleVectors(playerAngles, vecforward, NULL_VECTOR, NULL_VECTOR);
	vecforward[2] = 0.0;
	NormalizeVector(vecforward, vecforward);
	GetVectorCrossProduct(up, vecforward, right);
	float front = GetVectorDotProduct(vecDelta, vecforward);
	float side = GetVectorDotProduct(vecDelta, right);
	xpos = 360.0 * -front;
	ypos = 360.0 * -side;
	float flRotation = (ArcTangent2(xpos, ypos) + FLOAT_PI) * (57.29577951);
	float yawRadians = -flRotation * 0.017453293;
	xpos = (500 + (360.0 * Cosine(yawRadians))) / 1000.0;
	ypos = (500 - (360.0 * Sine(yawRadians))) / 1000.0;
	return;
}

public void KeyValuesToStringMap(const KeyValues kv, const StringMap stringmap, bool hide_top, int depth, char[] prefix)
{

	for (;;)
	{
		char section_name[128];

		kv.GetSectionName(section_name, sizeof(section_name));
		if (kv.GotoFirstSubKey(false))
		{
			char new_prefix[128];

			if ((depth == 0 && hide_top))
			{
				new_prefix = "";
			}
			else if ((prefix[0] == 0))
			{
				new_prefix = section_name;
			}
			else
			{
				FormatEx(new_prefix, sizeof(new_prefix), "%s.%s", prefix, section_name);
			}
			KeyValuesToStringMap(kv, stringmap, hide_top, depth + 1, new_prefix);
			kv.GoBack();
		}
		else 
		{
			if (kv.GetDataType(NULL_STRING) != KvData_None)
			{
				int keylen;

				char key[128];

				if (prefix[0] == 0)
				{
					key = section_name;
				}
				else 
				{
					FormatEx(key, sizeof(key), "%s.%s", prefix, section_name);
				}
				keylen = strlen(key);
				for (int i = 0; i < keylen; i++)
				{
					int bytes;

					bytes = IsCharMB(key[i]);
					if (bytes == 0)
					{
						key[i] = CharToLower(key[i]);
					}
					else 
					{
						i += (bytes - 1);
					}
				}
				char value[128];

				kv.GetString(NULL_STRING, value, sizeof(value), NULL_STRING);
				stringmap.SetValue(key, value);
			}
		}
		if (!kv.GotoNextKey(false))
		{
			break;
		}
	}
}

public int GetQueryRes(const DBResultSet dbr, float& GetQueryRes_param1)
{
	GetQueryRes_param1 = dbr.FetchFloat(1, null);
	return dbr.FetchInt(0, null);
}

public void SrcGoTmpFunc0(int client)
{
}

public Action SrcGoTmpFunc1(Handle timer, any data)
{
	return Plugin_Continue;
}

public int SrcGoTmpFunc2(int a, int b)
{
	return a + b;
}

public int SrcGoTmpFunc3(int a, int b, int& SrcGoTmpFunc3_param1)
{
	SrcGoTmpFunc3_param1 = a * b;
	return a + b;
}

public int SrcGoTmpFunc4(int a, int b)
{
	return a + b;
}

public bool SrcGoTmpFunc5(const char[] pattern, const ArrayList clients)
{
	bool non = StrContains(pattern, "!", false) != -1;
	for (int i = MAX_TF_PLAYERS; i > 0; i--)
	{
		if( IsClientValid(i) && clients.FindValue(i) == -1 ) {
			if( g_cvars.enabled.BoolValue && g_dnd.IsGameMaster(i) ) {
				if( !non ) {
					clients.Push(i);
				}
			} else if( non ) {
				clients.Push(i);
			}
		}
	}
	return true;
}