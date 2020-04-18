extern crate libc;

use libc::c_longlong;

use std::ffi::{CStr, CString};

#[derive(Debug)]
#[repr(C)]
struct GoString {
    a: *const libc::c_char,
    b: i64,
}

#[derive(Debug)]
#[repr(C)]
struct GoSlice {
    data: *mut libc::c_void,
    len: libc::c_longlong,
    cap: libc::c_longlong,
}

extern "C" {
    fn Add(a: c_longlong, b: c_longlong) -> c_longlong;
    fn AddArray(a: GoSlice, b: GoSlice, c: *mut GoSlice);
    fn AddString(a: GoString, b: GoString) -> *mut libc::c_char;
}

fn main() {
    let result = unsafe { Add(10i64, 12i64) };
    println!("{:?}", result);

    let mut a_arr = vec![20i64, 2i64].into_boxed_slice();
    let a_slice = GoSlice {
        data: a_arr.as_mut_ptr() as *mut libc::c_void,
        len: a_arr.len() as libc::c_longlong,
        cap: a_arr.len() as libc::c_longlong,
    };

    let mut b_arr = vec![138877474747i64, 2i64].into_boxed_slice();
    let b_slice = GoSlice {
        data: b_arr.as_mut_ptr() as *mut libc::c_void,
        len: b_arr.len() as libc::c_longlong,
        cap: b_arr.len() as libc::c_longlong,
    };
    println!("{:?}", a_slice);
    println!("{:?}", b_slice);

    let mut c_arr = vec![0i64, 0i64].into_boxed_slice();
    let mut c_slice = Box::new(GoSlice {
        data: c_arr.as_mut_ptr() as *mut libc::c_void,
        len: c_arr.len() as libc::c_longlong,
        cap: c_arr.len() as libc::c_longlong,
    });
    // 这段内存需要保持到返回为止
    std::mem::forget(c_arr);

    println!("begin to call");
    unsafe { AddArray(a_slice, b_slice, &mut *c_slice) };
    println!("end call");

    let c_arr = unsafe {
        std::vec::Vec::from_raw_parts(
            c_slice.data as *mut i64,
            c_slice.len as usize,
            c_slice.cap as usize,
        )
    };
    println!("{:?}", c_arr);
    std::mem::drop(c_arr);

    let c_path = CString::new("hello duanbing cstring").expect("CString::new failed");
    let ptr = c_path.as_ptr();
    let go_string = GoString {
        a: ptr,
        b: c_path.as_bytes().len() as i64,
    };
    let c_path = CString::new("hello duanbing cstring").expect("CString::new failed");
    let ptr = c_path.as_ptr();
    let go_string2 = GoString {
        a: ptr,
        b: c_path.as_bytes().len() as i64,
    };

    println!("go go ");
    let res = unsafe { AddString(go_string, go_string2) };
    let c_str = unsafe { CStr::from_ptr(res) };
    println!("res: {:?}", c_str.to_str().unwrap());
}
