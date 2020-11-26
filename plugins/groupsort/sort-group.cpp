#include "../../manticore/src/sphinxudf.h"

#include <string.h>
#include <iostream>
#include <vector>
#include <sstream>
#include <algorithm>

#ifdef _MSC_VER
#define snprintf _snprintf
#define DLLEXPORT __declspec(dllexport)
#else
#define DLLEXPORT extern "C" __attribute__((visibility("default")))
#endif
using std::vector;
using std::stringstream;
using std::sort;
/// UDF version control
/// gets called once when the library is loaded
DLLEXPORT int sort_group_ver() {
	return SPH_UDF_VERSION;
}
/// UDF initialization
/// gets called on every query, when query begins
/// args are filled with values for a particular query
DLLEXPORT int sort_group_init(SPH_UDF_INIT* init, SPH_UDF_ARGS* args, char* error_message) {
   // Mandatory testings
   if ((args->arg_count < 2) || (args->arg_count > 3) || (args->arg_types[0] != SPH_UDF_TYPE_STRING)) {
		snprintf(error_message, SPH_UDF_ERROR_LEN, "Impropeer value in args->arg_count or args->arg_types");
      return 1;
   }
   if ((args->str_lengths[0] == 0) || (args->str_lengths[1] == 0)) {
		snprintf(error_message, SPH_UDF_ERROR_LEN, "args->str_lengths[] can't be zero");
      return 1;
   }
	return 0;
}
/// UDF deinitialization
/// gets called on every query, when query ends
DLLEXPORT void sort_group_deinit(SPH_UDF_INIT* init) {
	// Nothing to destroy
}
/// UDF implementation
/// gets called for every row, unless optimized away
DLLEXPORT char* sort_group(SPH_UDF_INIT* init, SPH_UDF_ARGS* args, char* error_flag) {
   // This stupid code is neccessary because arg_values[] are not ASCIIZ strings
   char* sIDs, * sStatuses, * sRanks;
   sIDs = new char[args->str_lengths[0] + 1];
   //strncpy_s(sIDs, args->str_lengths[0] + 1, args->arg_values[0], args->str_lengths[0]);
   strncpy(sIDs, args->arg_values[0], args->str_lengths[0]);
   sStatuses = new char[args->str_lengths[1] + 1];
   //strncpy_s(sStatuses, args->str_lengths[1] + 1, args->arg_values[1], args->str_lengths[1]);
   strncpy(sStatuses, args->arg_values[1], args->str_lengths[1]);
   sRanks = new char[args->str_lengths[2] + 1];
   //strncpy_s(sRanks, args->str_lengths[2] + 1, args->arg_values[2], args->str_lengths[2]);
   strncpy(sRanks, args->arg_values[2], args->str_lengths[2]);
   // Init vectors
   int n, m, l;
   char ch;
   struct Value {
      int ID;
      int Status;
      int Rank;
   } v;
   std::vector<Value> values;
   stringstream streamIDs(sIDs);
   stringstream streamStatuses(sStatuses);
   stringstream streamRanks(sRanks);
   while (streamIDs >> n) {
      v.ID = n;
      streamStatuses >> m;
      v.Status = m;
      streamRanks >> l;
      v.Rank = l;
      values.push_back(v);
      streamIDs >> ch;
      streamStatuses >> ch;
      streamRanks >> ch;
      // Detect Statuses less than IDs
      if (!streamIDs.eof() && streamStatuses.eof()) {
         *error_flag = 1;
         return NULL;
      }
      // Detect Ranks less than IDs
      if (!streamIDs.eof() && streamRanks.eof() && (args->arg_count == 3)) {
         *error_flag = 1;
         return NULL;
      }
   }
   // Sort by Statuses
   std::sort(values.begin(), values.end(), [](Value a, Value b) {
      return a.Status < b.Status;
      });
   // Sort by Ranks inside Statuses
   int index = 0, index2 = 0;
   // Sort by Ranks inside Status=1
   std::vector<Value>::iterator it = std::find_if(values.begin(), values.end(), [](const Value& v) {
      return (v.Status == 2) ? true : false;
      });
   if (it != values.end()) {
      index = it - values.begin();
      std::sort(values.begin(), values.begin() + index, [](Value a, Value b) {
         return a.Rank < b.Rank;
         });
   }
   // Sort by Ranks inside Status=2
   it = std::find_if(values.begin() + index + 1, values.end(), [](const Value& v) {
      return (v.Status == 3) ? true : false;
      });
   if (it != values.end()) {
      index2 = it - values.begin();
      std::sort(values.begin() + index, values.begin() + index2, [](Value a, Value b) {
         return a.Rank < b.Rank;
         });
   }
   // Sort by Ranks inside Status=3
   if (it != values.end()) {
      std::sort(values.begin() + index2, values.end(), [](Value a, Value b) {
         return a.Rank < b.Rank;
         });
   }
   // Render strings
   std::stringstream().swap(streamIDs);
   std::stringstream().swap(streamStatuses);
   std::stringstream().swap(streamRanks);
   for (auto it = values.begin(); it != values.end(); ++it) {
      streamIDs << it->ID;
      streamStatuses << it->Status;
      if (args->arg_count == 3)
         streamRanks << it->Rank;
      if (std::next(it) != values.end()) {
         streamIDs << ",";
         streamStatuses << ",";
         if (args->arg_count == 3)
            streamRanks << ",";
      }
   }
   const std::string& tmpIDs = streamIDs.str();
   const char* pch = tmpIDs.c_str();
   // No way - we can't put \0 at the end!
   //memcpy_s(args->arg_values[0], args->str_lengths[0], pch, args->str_lengths[0]);
   memcpy(args->arg_values[0], pch, args->str_lengths[0]);
   const std::string& tmpStatuses = streamStatuses.str();
   pch = tmpStatuses.c_str();
   //memcpy_s(args->arg_values[1], args->str_lengths[1], pch, args->str_lengths[1]);
   memcpy(args->arg_values[1], pch, args->str_lengths[1]);
   if (args->arg_count == 3) {
      const std::string& tmpRanks = streamRanks.str();
      pch = tmpRanks.c_str();
      //memcpy_s(args->arg_values[2], args->str_lengths[2], pch, args->str_lengths[2]);
      memcpy(args->arg_values[2], pch, args->str_lengths[2]);
   }
   return NULL;
}