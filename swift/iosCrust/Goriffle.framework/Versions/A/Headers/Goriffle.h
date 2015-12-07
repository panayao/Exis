// Objective-C API for talking to github.com/exis-io/goriffle Go package.
//   gobind -lang=objc github.com/exis-io/goriffle
//
// File is generated by gobind. Do not edit.

#ifndef __GoGoriffle_H__
#define __GoGoriffle_H__

#include <Foundation/Foundation.h>

FOUNDATION_EXPORT NSString* const GoGoriffleACTION_SEPARATOR;
FOUNDATION_EXPORT NSString* const GoGoriffleDOMAIN_SEPARATOR;
FOUNDATION_EXPORT NSString* const GoGoriffleErrAuthorizationFailed;
FOUNDATION_EXPORT NSString* const GoGoriffleErrCloseRealm;
FOUNDATION_EXPORT NSString* const GoGoriffleErrDomainAlreadyExists;
FOUNDATION_EXPORT NSString* const GoGoriffleErrGoodbyeAndOut;
FOUNDATION_EXPORT NSString* const GoGoriffleErrInvalidArgument;
FOUNDATION_EXPORT NSString* const GoGoriffleErrInvalidUri;
FOUNDATION_EXPORT NSString* const GoGoriffleErrNoSuchDomain;
FOUNDATION_EXPORT NSString* const GoGoriffleErrNoSuchRealm;
FOUNDATION_EXPORT NSString* const GoGoriffleErrNoSuchRegistration;
FOUNDATION_EXPORT NSString* const GoGoriffleErrNoSuchRole;
FOUNDATION_EXPORT NSString* const GoGoriffleErrNoSuchSubscription;
FOUNDATION_EXPORT NSString* const GoGoriffleErrNotAuthorized;
FOUNDATION_EXPORT NSString* const GoGoriffleErrSystemShutdown;

FOUNDATION_EXPORT void GoGoriffleLog();

FOUNDATION_EXPORT NSString* GoGorifflePConnector(NSString* url, NSString* domain);

FOUNDATION_EXPORT NSData* GoGorifflePRecieve();

FOUNDATION_EXPORT NSData* GoGorifflePRegister(NSString* s);

FOUNDATION_EXPORT NSData* GoGorifflePSubscribe(NSString* s);

FOUNDATION_EXPORT void GoGorifflePYield(NSData* args);

#endif
