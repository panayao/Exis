package com.exis.riffle;

import com.exis.riffle.cumin.Cumin;
import com.exis.riffle.cumin.Handler;

import go.mantle.Mantle;

/**
 * Created by damouse on 2/15/16.
 *
 * Like the regular deferred, but with many more options when it comes to handling arguments
 */
public class CallDeferred extends Deferred {
    Mantle.Domain domain;

    public CallDeferred(AppDomain app, Mantle.Domain domain) {
        super(app);
        this.domain = domain;
    }

    public CallDeferred _then(Cumin.Wrapped fn, Object[] types) {
        _callback = fn;
        domain.CallExpects(cb.toString(), Utils.marshall(types));
        return this;
    }

    // Start Generic Shotgun

public  CallDeferred then( Handler.ZeroZero handler) {
    return _then(Cumin.cuminicate( handler), Cumin.representation());
}

public <A> CallDeferred then(Class<A> a,  Handler.OneZero<A> handler) {
    return _then(Cumin.cuminicate(a,  handler), Cumin.representation(a));
}

public <A, B> CallDeferred then(Class<A> a, Class<B> b,  Handler.TwoZero<A, B> handler) {
    return _then(Cumin.cuminicate(a, b,  handler), Cumin.representation(a, b));
}

public <A, B, C> CallDeferred then(Class<A> a, Class<B> b, Class<C> c,  Handler.ThreeZero<A, B, C> handler) {
    return _then(Cumin.cuminicate(a, b, c,  handler), Cumin.representation(a, b, c));
}

public <A, B, C, D> CallDeferred then(Class<A> a, Class<B> b, Class<C> c, Class<D> d,  Handler.FourZero<A, B, C, D> handler) {
    return _then(Cumin.cuminicate(a, b, c, d,  handler), Cumin.representation(a, b, c, d));
}

public <A, B, C, D, E> CallDeferred then(Class<A> a, Class<B> b, Class<C> c, Class<D> d, Class<E> e,  Handler.FiveZero<A, B, C, D, E> handler) {
    return _then(Cumin.cuminicate(a, b, c, d, e,  handler), Cumin.representation(a, b, c, d, e));
}

public <A, B, C, D, E, F> CallDeferred then(Class<A> a, Class<B> b, Class<C> c, Class<D> d, Class<E> e, Class<F> f,  Handler.SixZero<A, B, C, D, E, F> handler) {
    return _then(Cumin.cuminicate(a, b, c, d, e, f,  handler), Cumin.representation(a, b, c, d, e, f));
}
    // End Generic Shotgun
}
